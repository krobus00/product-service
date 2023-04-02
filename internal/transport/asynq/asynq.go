package asynq

import (
	"github.com/hibiken/asynq"
	"github.com/krobus00/product-service/internal/model"
	"github.com/sirupsen/logrus"
)

type Delivery struct {
	productUC model.ProductUsecase
	asynqMux  *asynq.ServeMux
}

func NewDelivery() *Delivery {
	return new(Delivery)
}

func (t *Delivery) InitRoutes() error {
	logrus.Info("register asynq handler")
	t.asynqMux.HandleFunc(model.TaskProductUpdateThumbnail, t.productUC.HandleUpdateThumbnailTask)

	return nil
}
