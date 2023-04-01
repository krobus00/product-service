//go:generate mockgen -destination=mock/mock_product_repository.go -package=mock github.com/krobus00/product-service/internal/model ProductRepository
//go:generate mockgen -destination=mock/mock_product_usecase.go -package=mock github.com/krobus00/product-service/internal/model ProductUsecase

package model

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	authPB "github.com/krobus00/auth-service/pb/auth"
	kit "github.com/krobus00/krokit"
	"github.com/krobus00/product-service/internal/utils"
	pb "github.com/krobus00/product-service/pb/product"
	storagePB "github.com/krobus00/storage-service/pb/storage"
	"gorm.io/gorm"
)

const (
	OSProductIndex              = "products"
	OSProductAnalyzer           = "my_analyzer"
	OSProductMinimumShouldMatch = "50%"

	ThumbnailType = string("IMAGE")
)

var (
	ProductSearchColumns = []string{"name", "description"}

	ErrProductNotFound         = errors.New("product not found")
	ErrThumbnailNotFound       = errors.New("thumbnail not found")
	ErrThumbnailTypeNotAllowed = errors.New("thumbnail type not allowed")
	ErrThumbnailNotAllowed     = errors.New("thumbnail not allowed")
)

type Product struct {
	ID          string `gorm:"primaryKey"`
	Name        string
	Description string
	Price       float64
	ThumbnailID string         // refer to object id
	OwnerID     string         // refer to user_id
	CreatedAt   time.Time      `gorm:"<-:create"` // read and create
	UpdatedAt   time.Time      `gorm:"<-"`        // allow read, create, and update
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

type Products []*Product

func (Product) TableName() string {
	return "products"
}

func NewProductCacheKey(id string) string {
	return fmt.Sprintf("products:id:%s", id)
}

func GetProductCacheKeys(id string) []string {
	return []string{
		NewProductCacheKey(id),
	}
}

func (m Products) ToProto() []*pb.Product {
	results := make([]*pb.Product, 0)
	for _, product := range m {
		if product == nil {
			continue
		}
		results = append(results, product.ToProto())
	}
	return results
}

func (m *Product) ToProto() *pb.Product {
	createdAt := m.CreatedAt.UTC().Format(time.RFC3339Nano)
	updatedAt := m.UpdatedAt.UTC().Format(time.RFC3339Nano)
	deletedAt := ""
	if m.DeletedAt.Valid {
		deletedAt = m.DeletedAt.Time.UTC().Format(time.RFC3339Nano)
	}
	return &pb.Product{
		Id:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       float32(m.Price),
		ThumbnailId: m.ThumbnailID,
		OwnerId:     m.OwnerID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		DeletedAt:   deletedAt,
	}
}

type DocProduct struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       float64        `json:"price"`
	OwnerID     string         `json:"owner_id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"deleted_at"`
}

func (m *DocProduct) GetID() string {
	return m.ID
}

func (m *Product) ToDoc() *DocProduct {
	deletedAt := gorm.DeletedAt{}
	if m.DeletedAt.Valid {
		deletedAt.Valid = true
		deletedAt.Time = m.DeletedAt.Time.UTC()
	}
	return &DocProduct{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		OwnerID:     m.OwnerID,
		CreatedAt:   m.CreatedAt.UTC(),
		UpdatedAt:   m.UpdatedAt.UTC(),
		DeletedAt:   deletedAt,
	}
}

type CreateProductPayload struct {
	ID          string
	Name        string
	Description string
	Price       float64
	ThumbnailID string // refer to object id
}

func NewCreateProductPayloadFromProto(message *pb.CreateProductRequest) *CreateProductPayload {
	return &CreateProductPayload{
		Name:        message.GetName(),
		Description: message.GetDescription(),
		Price:       float64(message.GetPrice()),
		ThumbnailID: message.GetThumbnailId(),
	}
}

func (m *CreateProductPayload) ToProduct(ownerID string) *Product {
	if m.ID == "" {
		m.ID = utils.GenerateUUID()
	}
	return &Product{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		Price:       m.Price,
		ThumbnailID: m.ThumbnailID,
		OwnerID:     ownerID,
	}
}

type UpdateProductPayload struct {
	ID          string
	Name        string
	Description string
	Price       float64
	ThumbnailID string // refer to object id
}

func NewUpdateProductPayloadFromProto(message *pb.UpdateProductRequest) *UpdateProductPayload {
	return &UpdateProductPayload{
		ID:          message.GetId(),
		Name:        message.GetName(),
		Description: message.GetDescription(),
		Price:       float64(message.GetPrice()),
		ThumbnailID: message.GetThumbnailId(),
	}
}

func (m *UpdateProductPayload) UpdateProduct(currentProduct *Product) *Product {
	currentProduct.Name = m.Name
	currentProduct.Description = m.Description
	currentProduct.Price = m.Price
	currentProduct.ThumbnailID = m.ThumbnailID
	return currentProduct
}

type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	Update(ctx context.Context, product *Product) error
	DeleteByID(ctx context.Context, id string) error
	FindPaginatedIDs(ctx context.Context, req *PaginationPayload) (ids []string, count int64, err error)
	FindOSPaginatedIDs(ctx context.Context, req *PaginationPayload) (ids []string, count int64, err error)

	// Resolver
	FindByID(ctx context.Context, id string) (*Product, error)

	// DI
	InjectDB(db *gorm.DB) error
	InjectRedisClient(client *redis.Client) error
	InjectOpensearchClient(client kit.OpensearchClient) error
}

type ProductUsecase interface {
	Create(ctx context.Context, payload *CreateProductPayload) (*Product, error)
	Update(ctx context.Context, payload *UpdateProductPayload) (*Product, error)
	Delete(ctx context.Context, id string) error
	FindPaginatedIDs(ctx context.Context, req *PaginationPayload) (*PaginationResponse, error)

	// Resolver
	FindByID(ctx context.Context, id string) (*Product, error)
	FindByIDs(ctx context.Context, ids []string) (Products, error)

	// DI
	InjectDB(db *gorm.DB) error
	InjectProductRepo(repo ProductRepository) error
	InjectAuthClient(client authPB.AuthServiceClient) error
	InjectStorageClient(client storagePB.StorageServiceClient) error
}
