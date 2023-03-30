package repository

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewProductRepository() model.ProductRepository {
	return new(productRepository)
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	logger := log.WithFields(log.Fields{
		"productID":   product.ID,
		"thumbnailID": product.ThumbnailID,
	})

	db := utils.GetTxFromContext(ctx, r.db)

	err := db.WithContext(ctx).Create(product).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetProductCacheKeys(product.ID))

	return nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	logger := log.WithFields(log.Fields{
		"productID": product.ID,
	})

	db := utils.GetTxFromContext(ctx, r.db)

	product.UpdatedAt = time.Now()
	err := db.WithContext(ctx).Updates(product).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetProductCacheKeys(product.ID))

	return nil
}

func (r *productRepository) DeleteByID(ctx context.Context, id string) error {
	logger := log.WithFields(log.Fields{
		"productID": id,
	})

	db := utils.GetTxFromContext(ctx, r.db)

	err := db.WithContext(ctx).Where("id = ?", id).Delete(&model.Product{}).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetProductCacheKeys(id))

	return nil
}

func (r *productRepository) FindPaginatedIDs(ctx context.Context, req *model.PaginationPayload) (ids []string, count int64, err error) {
	logger := log.WithFields(log.Fields{
		"search": req.Search,
		"sort":   req.Sort,
		"page":   req.Page,
		"limit":  req.Limit,
	})
	db := utils.GetTxFromContext(ctx, r.db)
	productIds := make([]string, 0)

	count, err = r.countPaginated(ctx, req)
	if err != nil {
		logger.Error(err.Error())
		return productIds, 0, err
	}

	err = db.WithContext(ctx).Scopes(
		WithPagination(req.Page, req.Limit),
		WithSearch(req.Search, model.ProductSearchColumns),
	).
		Select("id").
		Model(&model.Product{}).
		Pluck("id", &productIds).Error
	if err != nil {
		logger.Error(err.Error())
		return productIds, 0, err
	}

	return productIds, count, nil
}

func (r *productRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	logger := log.WithFields(log.Fields{
		"productID": id,
	})

	db := utils.GetTxFromContext(ctx, r.db)
	product := new(model.Product)
	cacheKey := model.NewProductCacheKey(id)

	cachedData, err := Get(ctx, r.redisClient, cacheKey)
	if err != nil {
		logger.Error(err.Error())
	}
	err = json.Unmarshal(cachedData, &product)
	if err == nil {
		return product, nil
	}

	product = new(model.Product)

	err = db.WithContext(ctx).Where("id = ?", id).First(product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = SetWithExpiry(ctx, r.redisClient, cacheKey, nil)
			if err != nil {
				logger.Error(err.Error())
			}
			return nil, nil
		}
		logger.Error(err.Error())
		return nil, err
	}

	err = SetWithExpiry(ctx, r.redisClient, cacheKey, product)
	if err != nil {
		logger.Error(err.Error())
	}

	return product, nil
}

func (r *productRepository) countPaginated(ctx context.Context, req *model.PaginationPayload) (count int64, err error) {
	logger := log.WithFields(log.Fields{
		"search": req.Search,
		"sort":   req.Sort,
		"page":   req.Page,
		"limit":  req.Limit,
	})
	db := utils.GetTxFromContext(ctx, r.db)
	err = db.WithContext(ctx).Scopes(
		WithSearch(req.Search, model.ProductSearchColumns),
	).
		Model(&model.Product{}).
		Select("id").
		Count(&count).Error

	if err != nil {
		logger.Error(err.Error())
		return 0, err
	}

	return count, nil
}
