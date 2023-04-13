package usecase

import (
	"context"
	"errors"
	"sync"

	"github.com/hibiken/asynq"
	authPB "github.com/krobus00/auth-service/pb/auth"
	"github.com/krobus00/product-service/internal/constant"
	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/utils"
	storagePB "github.com/krobus00/storage-service/pb/storage"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productUsecase struct {
	db            *gorm.DB
	productRepo   model.ProductRepository
	authClient    authPB.AuthServiceClient
	storageClient storagePB.StorageServiceClient
	jsClient      nats.JetStreamContext
	asynqClient   *asynq.Client
}

func NewProductUsecase() model.ProductUsecase {
	return new(productUsecase)
}

func (uc *productUsecase) Create(ctx context.Context, payload *model.CreateProductPayload) (*model.Product, error) {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	userID := getUserIDFromCtx(ctx)

	logger := logrus.WithFields(logrus.Fields{
		"userID": userID,
	})

	newProduct := payload.ToProduct(userID)

	err := uc.hasAccess(ctx, constant.ActionCreate, nil)
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
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	userID := getUserIDFromCtx(ctx)

	logger := logrus.WithFields(logrus.Fields{
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

	err = uc.hasAccess(ctx, constant.ActionUpdate, product)
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
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	userID := getUserIDFromCtx(ctx)

	logger := logrus.WithFields(logrus.Fields{
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

	if product.DeletedAt.Valid {
		return model.ErrProductAlreadyDeleted
	}

	err = uc.hasAccess(ctx, constant.ActionDelete, product)
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
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	var (
		ids        []string
		count      int64
		userID     = getUserIDFromCtx(ctx)
		dataSource = getDataSource(ctx)
	)

	logger := logrus.WithFields(logrus.Fields{
		"userID": userID,
		"search": req.Search,
		"sort":   req.Sort,
		"page":   req.Page,
		"limit":  req.Limit,
	})

	permission := []string{
		constant.PermissionProductAll,
	}

	if req.IncludeDeleted {
		permission = append(permission, constant.PermissionProductReadDeleted)
	} else {
		permission = append(permission, constant.PermissionProductRead, constant.PermissionProductReadOther)
	}

	err := hasAccess(ctx, uc.authClient, permission)
	if err != nil {
		logger.Error(err.Error())
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
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	userID := getUserIDFromCtx(ctx)

	logger := logrus.WithFields(logrus.Fields{
		"userID":    userID,
		"productID": id,
	})

	product, err := uc.productRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}
	if product == nil {
		return nil, model.ErrProductNotFound
	}

	err = uc.hasAccess(ctx, constant.ActionRead, product)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return product, nil
}

func (uc *productUsecase) FindByIDs(ctx context.Context, ids []string) (model.Products, error) {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	userID := getUserIDFromCtx(ctx)

	logger := logrus.WithFields(logrus.Fields{
		"userID":     userID,
		"productIDs": ids,
	})
	products := model.Products{}

	err := uc.hasAccess(ctx, constant.ActionRead, nil)
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
			logger := logger.WithFields(logrus.Fields{
				"userID":    userID,
				"productID": id,
			})
			product, err := uc.productRepo.FindByID(ctx, id)
			if err != nil {
				logger.Error(err.Error())
				return
			}
			if product != nil {
				productMapMu.Lock()
				productMap[id] = product
				productMapMu.Unlock()
			}
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

func (uc *productUsecase) hasAccess(ctx context.Context, action constant.ActionType, object *model.Product) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	userID := getUserIDFromCtx(ctx)

	permissions := []string{
		constant.PermissionProductAll,
	}

	switch action {
	case constant.ActionCreate:
		permissions = append(permissions, constant.PermissionProductCreate)
	case constant.ActionRead:
		if object == nil {
			permissions = append(permissions, constant.PermissionProductRead)
			break
		}
		if object.DeletedAt.Valid {
			permissions = append(permissions, constant.PermissionProductReadDeleted)
		}
		if !object.DeletedAt.Valid && object.OwnerID != userID {
			permissions = append(permissions, constant.PermissionProductReadOther)
		}
	case constant.ActionUpdate:
		if object.OwnerID == userID {
			permissions = append(permissions, constant.PermissionProductUpdate)
		} else {
			permissions = append(permissions, constant.PermissionProductModifyOther)
		}
	case constant.ActionDelete:
		if object.OwnerID == userID {
			permissions = append(permissions, constant.PermissionProductDelete)
		} else {
			permissions = append(permissions, constant.PermissionProductModifyOther)
		}
	default:
		return errors.New("invaid action")
	}

	err := hasAccess(ctx, uc.authClient, permissions)
	if err != nil {
		return err
	}

	return nil
}
