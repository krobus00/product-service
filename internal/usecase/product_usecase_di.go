package usecase

import (
	"errors"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/product-service/internal/model"
	storagePB "github.com/krobus00/storage-service/pb/storage"
	"gorm.io/gorm"
)

func (uc *productUsecase) InjectDB(db *gorm.DB) error {
	if db == nil {
		return errors.New("invalid db")
	}
	uc.db = db
	return nil
}

func (uc *productUsecase) InjectProductRepo(repo model.ProductRepository) error {
	if repo == nil {
		return errors.New("invalid product repository")
	}
	uc.productRepo = repo
	return nil
}

func (uc *productUsecase) InjectAuthClient(client authPB.AuthServiceClient) error {
	if client == nil {
		return errors.New("invalid auth client")
	}
	uc.authClient = client
	return nil
}

func (uc *productUsecase) InjectStorageClient(client storagePB.StorageServiceClient) error {
	if client == nil {
		return errors.New("invalid storage client")
	}
	uc.storageClient = client
	return nil
}
