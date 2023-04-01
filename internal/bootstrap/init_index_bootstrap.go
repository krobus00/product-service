package bootstrap

import (
	"context"

	"github.com/krobus00/product-service/internal/constant"
	"github.com/krobus00/product-service/internal/infrastructure"
	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/utils"
	"github.com/sirupsen/logrus"
)

func StartInitIndex() {
	osClient, err := infrastructure.NewOpensearchClient()
	utils.ContinueOrFatal(err)
	ctx := context.Background()

	logrus.Info("init index")
	resCreateIndices, err := osClient.CreateIndices(ctx, model.OSProductIndex, constant.InitIndex)
	if err != nil {
		logrus.Error(err.Error())
	}

	if resCreateIndices.IsError() {
		logrus.Error("error init index")
		logrus.Info(resCreateIndices)
	} else {
		logrus.Info("success init index")
	}

	resUpdateMapping, err := osClient.PutIndicesMapping(ctx, []string{model.OSProductIndex}, constant.IndexMapping)
	if err != nil {
		logrus.Error(err.Error())
	}

	if resUpdateMapping.IsError() {
		logrus.Error("error update mapping")
		logrus.Info(resUpdateMapping)
	} else {
		logrus.Info("success update mapping")
	}
}
