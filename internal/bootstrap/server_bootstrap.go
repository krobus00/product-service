package bootstrap

import (
	"context"
	"fmt"
	"net"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/product-service/internal/config"
	"github.com/krobus00/product-service/internal/infrastructure"
	"github.com/krobus00/product-service/internal/repository"
	grpcServer "github.com/krobus00/product-service/internal/transport/grpc"
	"github.com/krobus00/product-service/internal/usecase"
	"github.com/krobus00/product-service/internal/utils"
	pb "github.com/krobus00/product-service/pb/product"
	storagePB "github.com/krobus00/storage-service/pb/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	log "github.com/sirupsen/logrus"
)

func StartServer() {
	infrastructure.InitializeDBConn()

	// init infra
	db, err := infrastructure.DB.DB()
	utils.ContinueOrFatal(err)

	redisClient, err := infrastructure.NewRedisClient()
	utils.ContinueOrFatal(err)

	// init grpc client
	authConn, err := grpc.Dial(config.AuthGRPCHost(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.ContinueOrFatal(err)

	authClient := authPB.NewAuthServiceClient(authConn)

	storageConn, err := grpc.Dial(config.StorageGRPCHost(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	utils.ContinueOrFatal(err)
	storageClient := storagePB.NewStorageServiceClient(storageConn)

	// init repository
	productRepo := repository.NewProductRepository()
	err = productRepo.InjectDB(infrastructure.DB)
	utils.ContinueOrFatal(err)
	err = productRepo.InjectRedisClient(redisClient)
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

	// init grpc
	grpcDelivery := grpcServer.NewDelivery()
	err = grpcDelivery.InjectProductUsecase(productUsecase)
	utils.ContinueOrFatal(err)

	productGrpcServer := grpc.NewServer()

	pb.RegisterProductServiceServer(productGrpcServer, grpcDelivery)
	if config.Env() == "development" {
		reflection.Register(productGrpcServer)
	}
	lis, _ := net.Listen("tcp", ":"+config.GRPCport())

	go func() {
		_ = productGrpcServer.Serve(lis)
	}()
	log.Info(fmt.Sprintf("grpc server started on :%s", config.GRPCport()))

	wait := gracefulShutdown(context.Background(), config.GracefulShutdownTimeOut(), map[string]operation{
		"database connection": func(ctx context.Context) error {
			infrastructure.StopTickerCh <- true
			return db.Close()
		},
		"grpc": func(ctx context.Context) error {
			return lis.Close()
		},
	})

	<-wait
}
