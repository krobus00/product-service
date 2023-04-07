package usecase

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/utils"
	"github.com/sirupsen/logrus"
)

func (uc *productUsecase) HandleUpdateThumbnailTask(ctx context.Context, t *asynq.Task) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	payload := new(model.TaskUpdateThumbnailPayload)
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("json.Unmarshal failed: %v: %w", err, asynq.SkipRetry)
	}

	err := uc.productRepo.UpdateAllThumbnail(ctx, payload.OldObjectID, payload.NewObjectID)
	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	return nil
}
