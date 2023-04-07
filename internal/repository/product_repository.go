package repository

import (
	"context"
	"errors"
	"io"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/goccy/go-json"
	kit "github.com/krobus00/krokit"
	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type productRepository struct {
	db          *gorm.DB
	redisClient *redis.Client
	osClient    kit.OpensearchClient
}

func NewProductRepository() model.ProductRepository {
	return new(productRepository)
}

func (r *productRepository) Create(ctx context.Context, product *model.Product) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

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

	res, err := r.osClient.Index(ctx, model.OSProductIndex, product.ToDoc())
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer res.Body.Close()

	_ = DeleteByKeys(ctx, r.redisClient, model.GetProductCacheKeys(product.ID))

	return nil
}

func (r *productRepository) Update(ctx context.Context, product *model.Product) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

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

	res, err := r.osClient.Index(ctx, model.OSProductIndex, product.ToDoc())
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer res.Body.Close()

	_ = DeleteByKeys(ctx, r.redisClient, model.GetProductCacheKeys(product.ID))

	return nil
}

func (r *productRepository) DeleteByID(ctx context.Context, id string) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := log.WithFields(log.Fields{
		"productID": id,
	})

	db := utils.GetTxFromContext(ctx, r.db)

	product := new(model.Product)

	err := db.WithContext(ctx).Clauses(clause.Returning{}).
		Where("id = ?", id).Delete(product).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	res, err := r.osClient.Index(ctx, model.OSProductIndex, product.ToDoc())
	if err != nil {
		logger.Error(err.Error())
		return err
	}
	defer res.Body.Close()

	_ = DeleteByKeys(ctx, r.redisClient, model.GetProductCacheKeys(id))

	return nil
}

func (r *productRepository) FindPaginatedIDs(ctx context.Context, req *model.PaginationPayload) (ids []string, count int64, err error) {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

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

	err = db.WithContext(ctx).Unscoped().Scopes(
		WithPagination(req.Page, req.Limit),
		WithSearch(req.Search, model.ProductSearchColumns),
		WithSortBy(req.Sort),
		WithDeleted(req.IncludeDeleted),
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
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

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
	query := db.WithContext(ctx).Unscoped().Where("id = ?", id).First(product)

	err = query.Error
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
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := log.WithFields(log.Fields{
		"search": req.Search,
		"sort":   req.Sort,
		"page":   req.Page,
		"limit":  req.Limit,
	})
	db := utils.GetTxFromContext(ctx, r.db)
	err = db.WithContext(ctx).Unscoped().Scopes(
		WithSearch(req.Search, model.ProductSearchColumns),
		WithDeleted(req.IncludeDeleted),
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

func (r *productRepository) FindOSPaginatedIDs(ctx context.Context, req *model.PaginationPayload) (ids []string, count int64, err error) {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := log.WithFields(log.Fields{
		"search": req.Search,
		"sort":   req.Sort,
		"page":   req.Page,
		"limit":  req.Limit,
	})
	productIds := make([]string, 0)

	paginationRequest := &model.OSPaginationRequest{
		From:           int64((req.Page - 1) * req.Limit),
		Size:           int64(req.Limit),
		TrackTotalHits: true,
		Query: model.Query{
			Bool: model.Bool{
				Must: model.Must{
					MultiMatch: model.MultiMatch{
						Query:              req.Search,
						Analyzer:           model.OSProductAnalyzer,
						Fields:             model.ProductSearchColumns,
						MinimumShouldMatch: model.OSProductMinimumShouldMatch,
					},
				},
			},
		},
	}

	if !req.IncludeDeleted {
		paginationRequest.Query.Bool.Filter = append(paginationRequest.Query.Bool.Filter, model.Filter{
			Term: map[string]string{
				"deleted_at.keyword": "null",
			},
		})
	}

	paginationRequest.ParseSort(req)

	docData, err := json.Marshal(paginationRequest)
	if err != nil {
		logger.Error(err.Error())
		return productIds, count, err
	}

	body := strings.NewReader(string(docData))

	res, err := r.osClient.Search(ctx, []string{model.OSProductIndex}, body)
	if err != nil {
		logger.Error(err.Error())
		return productIds, count, err
	}
	defer res.Body.Close()

	osProducts := new(model.OSPaginationResponse[model.Product])

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		logger.Error(err.Error())
		return productIds, count, err
	}

	err = json.Unmarshal(bytes, &osProducts)
	if err != nil {
		logger.Error(err.Error())
		return productIds, count, err
	}

	for _, hit := range osProducts.Hits.Hits {
		productIds = append(productIds, hit.Source.ID)
	}

	return productIds, osProducts.GetCount(), nil
}

func (r *productRepository) UpdateAllThumbnail(ctx context.Context, oldThumbnailID string, newThumbnailID string) error {
	_, _, fn := utils.Trace()
	ctx, span := utils.NewSpan(ctx, fn)
	defer span.End()

	logger := log.WithFields(log.Fields{
		"oldThumbnailID": oldThumbnailID,
		"newThumbnailID": newThumbnailID,
	})

	db := utils.GetTxFromContext(ctx, r.db)

	err := db.WithContext(ctx).
		Model(&model.Product{}).
		Where("thumbnail_id = ?", oldThumbnailID).
		Updates(model.Product{ThumbnailID: newThumbnailID, UpdatedAt: time.Now()}).Error
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	_ = DeleteByKeys(ctx, r.redisClient, model.GetUpdateAllThumbnailCacheKeys())

	return nil
}
