package usecase

import (
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/krobus00/product-service/internal/config"
	"github.com/krobus00/product-service/internal/model"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
)

func (uc *productUsecase) CreateStream() error {
	stream, _ := uc.jsClient.StreamInfo(model.ProductStreamName)
	// stream not found, create it
	if stream == nil {
		logrus.Printf("Creating stream: %s\n", model.ProductStreamName)
		_, err := uc.jsClient.AddStream(&nats.StreamConfig{
			Name:     model.ProductStreamName,
			Subjects: []string{model.ProductStreamSubjects},
			MaxAge:   config.JetstreamMaxAge(),
			Storage:  nats.FileStorage,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (uc *productUsecase) ConsumeEvent() error {
	natsSubOpt := []nats.SubOpt{nats.ManualAck(), nats.Durable(config.DurableID())}

	logrus.Info(fmt.Sprintf("starting consume %s", model.ProductStreamSubjects))
	_, err := uc.jsClient.QueueSubscribe(model.ProductStreamSubjects, config.QueueGroup(), uc.consumeProductStream, natsSubOpt...)
	if err != nil {
		return err
	}

	return nil
}

func (uc *productUsecase) consumeProductStream(msg *nats.Msg) {
	err := msg.Ack()
	if err != nil {
		logrus.Error(fmt.Sprintf("unable to Ack: %v ", err))
		return
	}

	switch msg.Subject {
	case model.ProductThumbnailDeletedSubject:
		err = uc.handleProductThumbnailDeletedEvent(msg)
	default:
		logrus.Warn("unknown subject")
	}
	if err != nil {
		logrus.Error(err.Error())
		return
	}
}

func (uc *productUsecase) handleProductThumbnailDeletedEvent(msg *nats.Msg) error {
	msgPayload := new(model.JSDeleteObjectPayload)
	err := json.Unmarshal(msg.Data, &msgPayload)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	payload := &model.TaskUpdateThumbnailPayload{
		OldObjectID: msgPayload.ObjectID,
		NewObjectID: model.DefaultThumbnail,
	}

	taskPayload, err := json.Marshal(payload)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	_, err = uc.asynqClient.Enqueue(
		asynq.NewTask(model.TaskProductUpdateThumbnail, taskPayload),
		asynq.MaxRetry(config.AsynqRetry()),
		asynq.Retention(config.AsynqRetention()),
	)

	if err != nil {
		logrus.Error(err.Error())
		return err
	}
	return nil
}
