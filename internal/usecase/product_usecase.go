package usecase

import (
	"context"
	"sync"

	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/product-service/internal/constant"
	"github.com/krobus00/product-service/internal/model"
	storagePB "github.com/krobus00/storage-service/pb/storage"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productUsecase struct {
	db            *gorm.DB
	productRepo   model.ProductRepository
	authClient    authPB.AuthServiceClient
	storageClient storagePB.StorageServiceClient
}

func NewProductUsecase() model.ProductUsecase {
	return new(productUsecase)
}

func (uc *productUsecase) Create(ctx context.Context, payload *model.CreateProductPayload) (*model.Product, error) {
	userID := getUserIDFromCtx(ctx)

	logger := log.WithFields(log.Fields{
		"userID": userID,
	})

	newProduct := payload.ToProduct(userID)

	err := hasAccess(ctx, uc.authClient, []string{
		constant.PermissionProductAll,
		constant.PermissionProductCreate,
	})
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	object, err := uc.storageClient.GetObjectByID(ctx, &storagePB.GetObjectByIDRequest{
		UserId:   userID,
		ObjectId: payload.ThumbnailID,
	})
	if err != nil {
		logger.Error(err.Error())
		return nil, model.ErrThumbnailNotFound
	}

	if object.GetType() != model.ThumbnailType {
		return nil, model.ErrThumbnailTypeNotAllowed
	}

	if !object.GetIsPublic() {
		return nil, model.ErrThumbnailNotAllowed
	}

	err = uc.productRepo.Create(ctx, newProduct)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return newProduct, nil
}

func (uc *productUsecase) Update(ctx context.Context, payload *model.UpdateProductPayload) (*model.Product, error) {
	userID := getUserIDFromCtx(ctx)

	logger := log.WithFields(log.Fields{
		"userID": userID,
	})

	product, err := uc.productRepo.FindByID(ctx, payload.ID)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	if product == nil {
		return nil, model.ErrProductNotFound
	}

	err = uc.hasAccess(ctx, []string{
		constant.PermissionProductAll,
		constant.PermissionProductUpdate,
	}, product)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	object, err := uc.storageClient.GetObjectByID(ctx, &storagePB.GetObjectByIDRequest{
		UserId:   userID,
		ObjectId: payload.ThumbnailID,
	})
	if err != nil {
		logger.Error(err.Error())
		return nil, model.ErrThumbnailNotFound
	}

	if object.GetType() != model.ThumbnailType {
		return nil, model.ErrThumbnailTypeNotAllowed
	}

	if !object.GetIsPublic() {
		return nil, model.ErrThumbnailNotAllowed
	}

	product = payload.UpdateProduct(product)
	err = uc.productRepo.Update(ctx, product)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	return product, nil
}

func (uc *productUsecase) Delete(ctx context.Context, id string) error {
	userID := getUserIDFromCtx(ctx)

	logger := log.WithFields(log.Fields{
		"userID":    userID,
		"productID": id,
	})

	product, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	if product == nil {
		return model.ErrProductNotFound
	}

	err = uc.hasAccess(ctx, []string{
		constant.PermissionProductAll,
		constant.PermissionProductDelete,
	}, product)

	if err != nil {
		logger.Error(err.Error())
		return err
	}

	err = uc.productRepo.DeleteByID(ctx, product.ID)
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (uc *productUsecase) FindPaginatedIDs(ctx context.Context, req *model.PaginationPayload) (*model.PaginationResponse, error) {
	var (
		ids        = make([]string, 0)
		count      = int64(0)
		userID     = getUserIDFromCtx(ctx)
		dataSource = getDataSource(ctx)
	)

	logger := log.WithFields(log.Fields{
		"userID": userID,
		"search": req.Search,
		"sort":   req.Sort,
		"page":   req.Page,
		"limit":  req.Limit,
	})

	err := hasAccess(ctx, uc.authClient, []string{
		constant.PermissionProductAll,
		constant.PermissionProductRead,
	})
	if err != nil {
		return nil, err
	}

	req = req.Sanitize()
	switch dataSource {
	case constant.SourceDB:
		ids, count, err = uc.productRepo.FindPaginatedIDs(ctx, req)
	case constant.SourceOS:
		ids, count, err = uc.productRepo.FindOSPaginatedIDs(ctx, req)
	default:
		ids, count, err = uc.productRepo.FindOSPaginatedIDs(ctx, req)
	}
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	res := model.NewPaginationResponse(req).
		WithCount(count).
		WithItems(ids)

	return res.BuildResponse(), nil
}

func (uc *productUsecase) FindByID(ctx context.Context, id string) (*model.Product, error) {
	userID := getUserIDFromCtx(ctx)

	logger := log.WithFields(log.Fields{
		"userID":    userID,
		"productID": id,
	})

	err := hasAccess(ctx, uc.authClient, []string{
		constant.PermissionProductAll,
		constant.PermissionProductRead},
	)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	product, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	if product == nil {
		return nil, model.ErrProductNotFound
	}

	return product, nil
}

func (uc *productUsecase) FindByIDs(ctx context.Context, ids []string) (model.Products, error) {
	userID := getUserIDFromCtx(ctx)

	logger := log.WithFields(log.Fields{
		"userID":     userID,
		"productIDs": ids,
	})
	products := model.Products{}

	err := hasAccess(ctx, uc.authClient, []string{
		constant.PermissionProductAll,
		constant.PermissionProductRead,
	})
	if err != nil {
		logger.Error(err.Error())
		return products, err
	}

	productMap := map[string]*model.Product{}
	productMapMu := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(ids))
	for _, productID := range ids {
		go func(id string) {
			defer wg.Done()
			logger := logger.WithFields(log.Fields{
				"userID":    userID,
				"productID": id,
			})
			product, err := uc.productRepo.FindByID(ctx, id)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			productMapMu.Lock()
			productMap[id] = product
			productMapMu.Unlock()
		}(productID)
	}

	wg.Wait()

	for _, productID := range ids {
		if product, ok := productMap[productID]; ok {
			products = append(products, product)
		}
	}

	return products, nil
}

func (uc *productUsecase) hasAccess(ctx context.Context, permissions []string, object *model.Product) error {
	userID := getUserIDFromCtx(ctx)

	if object.OwnerID == userID {
		return nil
	}

	if object.OwnerID != userID {
		err := hasAccess(ctx, uc.authClient, []string{constant.PermissionProductModifyOther})
		if err != nil {
			return err
		}
		return nil
	}

	err := hasAccess(ctx, uc.authClient, permissions)
	if err != nil {
		return err
	}

	return nil
}
