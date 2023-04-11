package infrastructure

import (
	kit "github.com/krobus00/krokit"
	"github.com/krobus00/product-service/internal/config"
)

func NewOpensearchClient() (kit.OpensearchClient, error) {
	opensearchClient, err := kit.NewOpensearchClient(&kit.OSConfig{
		Addresses:          config.OpensearchHost(),
		InsecureSkipVerify: config.OpensearchInsecure(),
		Username:           config.OpensearchUsername(),
		Password:           config.OpensearchPassword(),
	})
	if err != nil {
		return nil, err
	}

	return opensearchClient, err
}
