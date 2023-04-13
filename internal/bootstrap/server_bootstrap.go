package bootstrap

import (
	"context"
	"fmt"
	"net"
	"net/http"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/product-service/internal/config"
	"github.com/krobus00/product-service/internal/infrastructure"
	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/repository"
	grpcServer "github.com/krobus00/product-service/internal/transport/grpc"
	"github.com/krobus00/product-service/internal/usecase"
	"github.com/krobus00/product-service/internal/utils"
	pb "github.com/krobus00/product-service/pb/product"
	storagePB "github.com/krobus00/storage-service/pb/storage"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func StartServer() {
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

	// init stream
	publisherUsecase := []model.PublisherUsecase{
		productUsecase,
	}

	for _, uc := range publisherUsecase {
		err = uc.CreateStream()
		utils.ContinueOrFatal(err)
	}

	// init grpc
	grpcDelivery := grpcServer.NewDelivery()
	err = grpcDelivery.InjectProductUsecase(productUsecase)
	utils.ContinueOrFatal(err)

	productGrpcServer := grpc.NewServer()

	pb.RegisterProductServiceServer(productGrpcServer, grpcDelivery)
	if config.Env() == "development" {
		reflection.Register(productGrpcServer)
	}
	lis, _ := net.Listen("tcp", ":"+config.PortGRPC())

	go func() {
		_ = productGrpcServer.Serve(lis)
	}()
	logrus.Info(fmt.Sprintf("grpc server started on :%s", config.PortGRPC()))

	http.Handle("/metrics", promhttp.Handler())

	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%s", config.PortMetrics()), nil)
	}()
	logrus.Info(fmt.Sprintf("metrics server started on :%s", config.PortMetrics()))

	wait := gracefulShutdown(context.Background(), config.GracefulShutdownTimeOut(), map[string]operation{
		"database connection": func(ctx context.Context) error {
			infrastructure.StopTickerCh <- true
			return db.Close()
		},
		"grpc": func(ctx context.Context) error {
			return lis.Close()
		},
		"nats connection": func(ctx context.Context) error {
			return nc.Drain()
		},
		"trace provider": func(ctx context.Context) error {
			if config.DisableTracing() {
				return nil
			}
			return tp.Shutdown(ctx)
		},
	})

	<-wait
}
