package repository

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/alicebob/miniredis/v2"
	"github.com/golang/mock/gomock"
	kitMock "github.com/krobus00/krokit/mock"
	"github.com/krobus00/product-service/internal/infrastructure"
	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/utils"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func newProductRepoMock(t *testing.T) (model.ProductRepository, sqlmock.Sqlmock, *miniredis.Miniredis) {
	db, sqlMock := utils.NewDBMock()
	miniRedis := miniredis.RunT(t)
	viper.Set("redis.cache_host", fmt.Sprintf("redis://%s", miniRedis.Addr()))
	redisClient, err := infrastructure.NewRedisClient()
	utils.ContinueOrFatal(err)
	productRepo := NewProductRepository()
	err = productRepo.InjectDB(db)
	utils.ContinueOrFatal(err)
	err = productRepo.InjectRedisClient(redisClient)
	utils.ContinueOrFatal(err)

	return productRepo, sqlMock, miniRedis
}

func Test_productRepository_Create(t *testing.T) {
	productID := utils.GenerateUUID()
	type mockIndex struct {
		err error
	}
	type args struct {
		product *model.Product
	}
	tests := []struct {
		name      string
		args      args
		mockIndex *mockIndex
		mockErr   error
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				product: &model.Product{
					ID:          productID,
					Name:        "new product",
					Description: "product description",
					Price:       17.17,
					ThumbnailID: utils.GenerateUUID(),
					OwnerID:     utils.GenerateUUID(),
				},
			},
			mockIndex: &mockIndex{
				err: nil,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				product: &model.Product{
					ID:          productID,
					Name:        "new product",
					Description: "product description",
					Price:       17.17,
					ThumbnailID: utils.GenerateUUID(),
					OwnerID:     utils.GenerateUUID(),
				},
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r, dbMock, _ := newProductRepoMock(t)
			osClient := kitMock.NewMockOpensearchClient(ctrl)

			err := r.InjectOpensearchClient(osClient)
			utils.ContinueOrFatal(err)

			dbMock.ExpectBegin()
			dbMock.ExpectExec("INSERT INTO \"products\"").
				WithArgs(productID, tt.args.product.Name, tt.args.product.Description, tt.args.product.Price, tt.args.product.ThumbnailID, tt.args.product.OwnerID, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.mockErr)

			if tt.mockIndex != nil {
				stringReader := strings.NewReader("test")
				stringReadCloser := ioutil.NopCloser(stringReader)
				osClient.EXPECT().Index(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&opensearchapi.Response{
					Body: stringReadCloser,
				}, tt.mockIndex.err)
			}

			if tt.wantErr {
				dbMock.ExpectRollback()
			} else {
				dbMock.ExpectCommit()
			}

			if err := r.Create(context.TODO(), tt.args.product); (err != nil) != tt.wantErr {
				t.Errorf("productRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_productRepository_Update(t *testing.T) {
	productID := utils.GenerateUUID()
	type mockIndex struct {
		err error
	}
	type args struct {
		product *model.Product
	}
	tests := []struct {
		name      string
		args      args
		mockIndex *mockIndex
		mockErr   error
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				product: &model.Product{
					ID:          productID,
					Name:        "new product",
					Description: "product description",
					Price:       17.17,
					ThumbnailID: utils.GenerateUUID(),
					OwnerID:     utils.GenerateUUID(),
				},
			},
			mockIndex: &mockIndex{
				err: nil,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				product: &model.Product{
					ID:          productID,
					Name:        "new product",
					Description: "product description",
					Price:       17.17,
					ThumbnailID: utils.GenerateUUID(),
					OwnerID:     utils.GenerateUUID(),
				},
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r, dbMock, _ := newProductRepoMock(t)
			osClient := kitMock.NewMockOpensearchClient(ctrl)

			err := r.InjectOpensearchClient(osClient)
			utils.ContinueOrFatal(err)

			dbMock.ExpectBegin()
			dbMock.ExpectExec("UPDATE \"products\"").
				WithArgs(tt.args.product.Name, tt.args.product.Description, tt.args.product.Price, tt.args.product.ThumbnailID, tt.args.product.OwnerID, sqlmock.AnyArg(), productID).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.mockErr)

			if tt.mockIndex != nil {
				stringReader := strings.NewReader("test")
				stringReadCloser := ioutil.NopCloser(stringReader)
				osClient.EXPECT().Index(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&opensearchapi.Response{
					Body: stringReadCloser,
				}, tt.mockIndex.err)
			}

			if tt.wantErr {
				dbMock.ExpectRollback()
			} else {
				dbMock.ExpectCommit()
			}
			if err := r.Update(context.TODO(), tt.args.product); (err != nil) != tt.wantErr {
				t.Errorf("productRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_productRepository_DeleteByID(t *testing.T) {
	productID := utils.GenerateUUID()
	type mockIndex struct {
		err error
	}
	type args struct {
		id string
	}
	tests := []struct {
		name      string
		args      args
		mockIndex *mockIndex
		mockErr   error
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				id: productID,
			},
			mockIndex: &mockIndex{
				err: nil,
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				id: productID,
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r, dbMock, _ := newProductRepoMock(t)
			osClient := kitMock.NewMockOpensearchClient(ctrl)

			err := r.InjectOpensearchClient(osClient)
			utils.ContinueOrFatal(err)

			dbMock.ExpectBegin()
			row := sqlmock.NewRows([]string{"id"})

			row.AddRow(tt.args.id)

			dbMock.ExpectQuery("UPDATE \"products\"").
				WithArgs(sqlmock.AnyArg(), tt.args.id).
				WillReturnRows(row).
				WillReturnError(tt.mockErr)

			if tt.mockIndex != nil {
				stringReader := strings.NewReader("test")
				stringReadCloser := ioutil.NopCloser(stringReader)
				osClient.EXPECT().Index(gomock.Any(), gomock.Any(), gomock.Any()).Times(1).Return(&opensearchapi.Response{
					Body: stringReadCloser,
				}, tt.mockIndex.err)
			}

			if tt.wantErr {
				dbMock.ExpectRollback()
			} else {
				dbMock.ExpectCommit()
			}
			if err := r.DeleteByID(context.TODO(), tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("productRepository.DeleteByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_productRepository_FindPaginatedIDs(t *testing.T) {
	productIds := []string{utils.GenerateUUID(), utils.GenerateUUID()}
	type args struct {
		req *model.PaginationPayload
	}
	type mockCount struct {
		count int64
		err   error
	}
	type mockSelect struct {
		ids []string
		err error
	}
	tests := []struct {
		name       string
		args       args
		mockCount  *mockCount
		mockSelect *mockSelect
		wantIds    []string
		wantCount  int64
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				req: &model.PaginationPayload{
					Search: "search something",
					Sort:   []string{},
					Limit:  10,
					Page:   1,
				},
			},
			mockCount: &mockCount{
				count: int64(len(productIds)),
				err:   nil,
			},
			mockSelect: &mockSelect{
				ids: productIds,
				err: nil,
			},
			wantIds:   productIds,
			wantCount: int64(len(productIds)),
			wantErr:   false,
		},
		{
			name: "count error",
			args: args{
				req: &model.PaginationPayload{
					Search: "",
					Sort:   []string{},
					Limit:  10,
					Page:   1,
				},
			},
			mockCount: &mockCount{
				count: 0,
				err:   errors.New("count error"),
			},

			wantIds:   []string{},
			wantCount: 0,
			wantErr:   true,
		},
		{
			name: "select error",
			args: args{
				req: &model.PaginationPayload{
					Search: "search something",
					Sort:   []string{},
					Limit:  10,
					Page:   1,
				},
			},
			mockCount: &mockCount{
				count: int64(len(productIds)),
				err:   nil,
			},
			mockSelect: &mockSelect{
				ids: []string{},
				err: errors.New("select error"),
			},
			wantIds:   []string{},
			wantCount: 0,
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, dbMock, _ := newProductRepoMock(t)

			if tt.mockCount != nil {
				row := sqlmock.NewRows([]string{"count"}).
					AddRow(tt.mockCount.count)
				dbMock.ExpectQuery("^SELECT COUNT.+ FROM \"products\"").
					// WithArgs().
					WillReturnRows(row).
					WillReturnError(tt.mockCount.err)
			}
			if tt.mockSelect != nil {
				row := sqlmock.NewRows([]string{"id"})
				for _, id := range tt.mockSelect.ids {
					row.AddRow(id)
				}
				dbMock.ExpectQuery("^SELECT .+ FROM \"products\"").
					WillReturnRows(row).
					WillReturnError(tt.mockSelect.err)
			}

			gotIds, gotCount, err := r.FindPaginatedIDs(context.TODO(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("productRepository.FindPaginatedIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(gotIds, tt.wantIds) {
				t.Errorf("productRepository.FindPaginatedIDs() gotIds = %v, want %v", gotIds, tt.wantIds)
			}
			if gotCount != tt.wantCount {
				t.Errorf("productRepository.FindPaginatedIDs() gotCount = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}

func Test_productRepository_FindByID(t *testing.T) {
	productID := utils.GenerateUUID()
	type args struct {
		id string
	}
	type mockSelect struct {
		product *model.Product
		err     error
	}
	tests := []struct {
		name       string
		args       args
		mockSelect *mockSelect
		want       *model.Product
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				id: productID,
			},
			mockSelect: &mockSelect{
				product: &model.Product{
					ID:          productID,
					Name:        "product-1",
					Description: "product description",
					Price:       17.17,
					ThumbnailID: "thumbnail-uuid",
					OwnerID:     "owner-uuid",
				},

				err: nil,
			},
			want: &model.Product{
				ID:          productID,
				Name:        "product-1",
				Description: "product description",
				Price:       17.17,
				ThumbnailID: "thumbnail-uuid",
				OwnerID:     "owner-uuid",
			},
			wantErr: false,
		},
		{
			name: "error record not found",
			args: args{
				id: productID,
			},
			mockSelect: &mockSelect{
				product: nil,
				err:     gorm.ErrRecordNotFound,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				id: productID,
			},
			mockSelect: &mockSelect{
				product: nil,
				err:     errors.New("db error"),
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, dbMock, _ := newProductRepoMock(t)
			if tt.mockSelect != nil {
				row := sqlmock.NewRows([]string{"id", "name", "description", "price", "thumbnail_id", "owner_id", "created_at", "updated_at", "deleted_at"})
				if tt.mockSelect.product != nil {
					product := tt.mockSelect.product
					row.AddRow(product.ID, product.Name, product.Description, product.Price, product.ThumbnailID, product.OwnerID, product.CreatedAt, product.UpdatedAt, product.DeletedAt)
				}

				dbMock.ExpectQuery("^SELECT .+ FROM \"products\"").
					WithArgs(productID).
					WillReturnRows(row).
					WillReturnError(tt.mockSelect.err)
			}
			got, err := r.FindByID(context.TODO(), tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("productRepository.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productRepository.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_productRepository_UpdateAllThumbnail(t *testing.T) {
	type args struct {
		oldThumbnailID string
		newThumbnailID string
	}
	tests := []struct {
		name    string
		args    args
		mockErr error
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				oldThumbnailID: "123",
				newThumbnailID: "321",
			},
			mockErr: nil,
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				oldThumbnailID: "123",
				newThumbnailID: "321",
			},
			mockErr: errors.New("db error"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r, dbMock, _ := newProductRepoMock(t)

			dbMock.ExpectBegin()
			dbMock.ExpectExec("UPDATE \"products\"").
				WithArgs(tt.args.newThumbnailID, sqlmock.AnyArg(), tt.args.oldThumbnailID).
				WillReturnResult(sqlmock.NewResult(1, 1)).
				WillReturnError(tt.mockErr)

			if tt.wantErr {
				dbMock.ExpectRollback()
			} else {
				dbMock.ExpectCommit()
			}

			if err := r.UpdateAllThumbnail(context.TODO(), tt.args.oldThumbnailID, tt.args.newThumbnailID); (err != nil) != tt.wantErr {
				t.Errorf("productRepository.UpdateAllThumbnail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_productRepository_FindOSPaginatedIDs(t *testing.T) {
	type args struct {
		req *model.PaginationPayload
	}
	type osMock struct {
		resp string
		err  error
	}
	tests := []struct {
		name      string
		args      args
		osMock    *osMock
		wantIds   []string
		wantCount int64
		wantErr   bool
	}{
		{
			name: "success",
			args: args{
				req: &model.PaginationPayload{
					Search: "sample product",
					Sort:   []string{},
					Limit:  10,
					Page:   1,
				},
			},
			osMock: &osMock{
				resp: `{
          "took": 33,
          "timed_out": false,
          "_shards": {
            "total": 1,
            "successful": 1,
            "skipped": 0,
            "failed": 0
          },
          "hits": {
            "total": {
              "value": 1,
              "relation": "eq"
            },
            "max_score": 45.986057,
            "hits": [
              {
                "_index": "products",
                "_id": "1fc34a8d-3d77-4f40-acbc-789c06b4fa5d",
                "_score": 45.986057,
                "_source": {
                  "id": "1fc34a8d-3d77-4f40-acbc-789c06b4fa5d",
                  "name": "product-",
                  "description": "description-1680269152107",
                  "price": 5222789120,
                  "owner_id": "cd9614c8-112a-4374-9737-eb62cc5d6aef",
                  "created_at": "2023-03-31T13:25:52.448676961Z",
                  "updated_at": "2023-03-31T13:25:52.448676961Z",
                  "deleted_at": null
                }
              }
            ]
          }
        }`,
				err: nil,
			},
			wantIds:   []string{"1fc34a8d-3d77-4f40-acbc-789c06b4fa5d"},
			wantCount: 1,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			r, _, _ := newProductRepoMock(t)
			osClient := kitMock.NewMockOpensearchClient(ctrl)

			err := r.InjectOpensearchClient(osClient)
			utils.ContinueOrFatal(err)

			if tt.osMock != nil {
				body := io.NopCloser(strings.NewReader(tt.osMock.resp))
				osClient.EXPECT().Search(gomock.Any(), gomock.Any(), gomock.Any()).
					Times(1).
					Return(&opensearchapi.Response{
						StatusCode: 200,
						Body:       body,
					}, tt.osMock.err)
			}

			gotIds, gotCount, err := r.FindOSPaginatedIDs(context.TODO(), tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("productRepository.FindOSPaginatedIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotIds, tt.wantIds) {
				t.Errorf("productRepository.FindOSPaginatedIDs() gotIds = %v, want %v", gotIds, tt.wantIds)
			}
			if gotCount != tt.wantCount {
				t.Errorf("productRepository.FindOSPaginatedIDs() gotCount = %v, want %v", gotCount, tt.wantCount)
			}
		})
	}
}
