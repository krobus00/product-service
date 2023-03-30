package grpc

import (
	"errors"

	"github.com/krobus00/product-service/internal/model"
)

func (t *Delivery) InjectProductUsecase(uc model.ProductUsecase) error {
	if uc == nil {
		return errors.New("invalid product usecase")
	}
	t.productUC = uc
	return nil
}
