package bootstrap

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hibiken/asynq"
	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/product-service/internal/config"
	"github.com/krobus00/product-service/internal/infrastructure"
	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/repository"
	asynqTransport "github.com/krobus00/product-service/internal/transport/asynq"
	"github.com/krobus00/product-service/internal/usecase"
	"github.com/krobus00/product-service/internal/utils"
	storagePB "github.com/krobus00/storage-service/pb/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartWorker() {
	infrastructure.InitializeDBConn()

	// init infra
	db, err := infrastructure.DB.DB()
	utils.ContinueOrFatal(err)

	redisClient, err := infrastructure.NewRedisClient()
	utils.ContinueOrFatal(err)

	osClient, err := infrastructure.NewOpensearchClient()
	utils.ContinueOrFatal(err)

	nc, js, err := infrastructure.NewJetstreamClient()
	utils.ContinueOrFatal(err)

	asynqClient, err := infrastructure.NewAsynqClient()
	utils.ContinueOrFatal(err)
	asynqServer, err := infrastructure.NewAsynqServer()
	utils.ContinueOrFatal(err)

	tp, err := infrastructure.JaegerTraceProvider()
	utils.ContinueOrFatal(err)

	// init grpc client
	grpcOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}

	authConn, err := grpc.Dial(config.AuthGRPCHost(), grpcOpts...)
	utils.ContinueOrFatal(err)
	authClient := authPB.NewAuthServiceClient(authConn)

	storageConn, err := grpc.Dial(config.StorageGRPCHost(), grpcOpts...)
	utils.ContinueOrFatal(err)
	storageClient := storagePB.NewStorageServiceClient(storageConn)

	// init repository
	productRepo := repository.NewProductRepository()
	err = productRepo.InjectDB(infrastructure.DB)
	utils.ContinueOrFatal(err)
	err = productRepo.InjectRedisClient(redisClient)
	utils.ContinueOrFatal(err)
	err = productRepo.InjectOpensearchClient(osClient)
	utils.ContinueOrFatal(err)

	// init usecase
	productUsecase := usecase.NewProductUsecase()
	err = productUsecase.InjectProductRepo(productRepo)
	utils.ContinueOrFatal(err)
	err = productUsecase.InjectDB(infrastructure.DB)
	utils.ContinueOrFatal(err)
	err = productUsecase.InjectAuthClient(authClient)
	utils.ContinueOrFatal(err)
	err = productUsecase.InjectStorageClient(storageClient)
	utils.ContinueOrFatal(err)
	err = productUsecase.InjectJetstreamClient(js)
	utils.ContinueOrFatal(err)
	err = productUsecase.InjectAsynqClient(asynqClient)
	utils.ContinueOrFatal(err)

	// init stream
	consumerUsecase := []model.ConsumerUsecase{
		productUsecase,
	}

	for _, uc := range consumerUsecase {
		go func(uc model.ConsumerUsecase) {
			err = uc.ConsumeEvent()
			utils.ContinueOrFatal(err)
		}(uc)
	}

	// init asnyq
	mux := asynq.NewServeMux()
	asynqDelivery := asynqTransport.NewDelivery()
	err = asynqDelivery.InjectProductUsecase(productUsecase)
	utils.ContinueOrFatal(err)
	err = asynqDelivery.InjectAsynqMux(mux)
	utils.ContinueOrFatal(err)

	err = asynqDelivery.InitRoutes()
	utils.ContinueOrFatal(err)

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%s", config.PortMetrics()), nil)
	}()
	logrus.Info(fmt.Sprintf("metrics server started on :%s", config.PortMetrics()))

	if err := asynqServer.Run(mux); err != nil {
		logrus.Fatalf("could not run asynq server: %v", err)
	}

	wait := gracefulShutdown(context.Background(), config.GracefulShutdownTimeOut(), map[string]operation{
		"database connection": func(ctx context.Context) error {
			infrastructure.StopTickerCh <- true
			return db.Close()
		},
		"nats connection": func(ctx context.Context) error {
			return nc.Drain()
		},
		"asynq client connection": func(ctx context.Context) error {
			return asynqClient.Close()
		},
		"asynq server connection": func(ctx context.Context) error {
			asynqServer.Shutdown()
			return asynqClient.Close()
		},
		"trace provider": func(ctx context.Context) error {
			return tp.Shutdown(ctx)
		},
	})

	<-wait
}
