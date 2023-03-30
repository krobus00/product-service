package grpc

import (
	"github.com/krobus00/product-service/internal/model"
	pb "github.com/krobus00/product-service/pb/product"
)

type Delivery struct {
	productUC model.ProductUsecase
	pb.UnsafeProductServiceServer
}

func NewDelivery() *Delivery {
	return new(Delivery)
}
