package asynq

import (
	"errors"

	"github.com/hibiken/asynq"
	"github.com/krobus00/product-service/internal/model"
)

func (t *Delivery) InjectProductUsecase(uc model.ProductUsecase) error {
	if uc == nil {
		return errors.New("invalid product usecase")
	}
	t.productUC = uc
	return nil
}

func (t *Delivery) InjectAsynqMux(mux *asynq.ServeMux) error {
	if mux == nil {
		return errors.New("invald asynq mux")
	}
	t.asynqMux = mux
	return nil
}
