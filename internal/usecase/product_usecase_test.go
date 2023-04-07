package usecase

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	authMock "github.com/krobus00/auth-service/pb/auth/mock"
	"github.com/krobus00/product-service/internal/constant"
	"github.com/krobus00/product-service/internal/model"
	"github.com/krobus00/product-service/internal/model/mock"
	"github.com/krobus00/product-service/internal/utils"
	storagePB "github.com/krobus00/storage-service/pb/storage"
	storageMock "github.com/krobus00/storage-service/pb/storage/mock"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func Test_productUsecase_Create(t *testing.T) {
	userID := utils.GenerateUUID()
	productID := utils.GenerateUUID()
	thumbnailID := utils.GenerateUUID()
	type args struct {
		payload *model.CreateProductPayload
	}
	type mockGetObjectByID struct {
		res *storagePB.Object
		err error
	}
	type mockCreate struct {
		err error
	}
	type mockAuth struct {
		hasAccess bool
		err       error
	}

	tests := []struct {
		name              string
		args              args
		userID            string
		mockGetObjectByID *mockGetObjectByID
		mockCreate        *mockCreate
		mockAuth          *mockAuth
		want              *model.Product
		wantErr           bool
	}{
		{
			name: "success",
			args: args{
				payload: &model.CreateProductPayload{
					ID:          productID,
					Name:        "new product",
					Description: "product description",
					Price:       17.17,
					ThumbnailID: thumbnailID,
				},
			},
			userID: userID,
			mockGetObjectByID: &mockGetObjectByID{
				res: &storagePB.Object{
					Id:         thumbnailID,
					FileName:   "test.png",
					Type:       model.ThumbnailType,
					SignedUrl:  "url",
					IsPublic:   true,
					UploadedBy: userID,
				},
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			mockCreate: &mockCreate{
				err: nil,
			},
			want: &model.Product{
				ID:          productID,
				Name:        "new product",
				Description: "product description",
				Price:       17.17,
				ThumbnailID: thumbnailID,
				OwnerID:     userID,
			},
			wantErr: false,
		},
		{
			name: "db error",
			args: args{
				payload: &model.CreateProductPayload{
					ID:          productID,
					Name:        "new product",
					Description: "product description",
					Price:       17.17,
					ThumbnailID: thumbnailID,
				},
			},
			userID: userID,
			mockGetObjectByID: &mockGetObjectByID{
				res: &storagePB.Object{
					Id:         thumbnailID,
					FileName:   "test.png",
					Type:       model.ThumbnailType,
					SignedUrl:  "url",
					IsPublic:   true,
					UploadedBy: userID,
				},
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			mockCreate: &mockCreate{
				err: errors.New("db error"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "permission denied",
			args: args{
				payload: &model.CreateProductPayload{
					ID:          productID,
					Name:        "new product",
					Description: "product description",
					Price:       17.17,
					ThumbnailID: thumbnailID,
				},
			},
			userID: userID,
			mockAuth: &mockAuth{
				hasAccess: false,
				err:       nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			ctx = context.WithValue(ctx, constant.KeyUserIDCtx, tt.userID)

			uc := NewProductUsecase()
			db, _ := utils.NewDBMock()
			err := uc.InjectDB(db)
			utils.ContinueOrFatal(err)
			mockAuthClient := authMock.NewMockAuthServiceClient(ctrl)
			err = uc.InjectAuthClient(mockAuthClient)
			utils.ContinueOrFatal(err)
			mockStorageClient := storageMock.NewMockStorageServiceClient(ctrl)
			err = uc.InjectStorageClient(mockStorageClient)
			utils.ContinueOrFatal(err)
			mockProductRepo := mock.NewMockProductRepository(ctrl)
			err = uc.InjectProductRepo(mockProductRepo)
			utils.ContinueOrFatal(err)

			if tt.mockAuth != nil {
				mockAuthClient.EXPECT().HasAccess(gomock.Any(), gomock.Any()).Times(1).Return(&wrapperspb.BoolValue{
					Value: tt.mockAuth.hasAccess,
				}, tt.mockAuth.err)
			}

			if tt.mockCreate != nil {
				product := tt.args.payload.ToProduct(userID)
				mockProductRepo.EXPECT().Create(gomock.Any(), product).Times(1).Return(tt.mockCreate.err)
			}

			if tt.mockGetObjectByID != nil {
				mockStorageClient.EXPECT().GetObjectByID(gomock.Any(), &storagePB.GetObjectByIDRequest{
					UserId:   tt.userID,
					ObjectId: tt.args.payload.ThumbnailID,
				}).Times(1).Return(tt.mockGetObjectByID.res, tt.mockGetObjectByID.err)
			}

			got, err := uc.Create(ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("productUsecase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productUsecase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_productUsecase_Update(t *testing.T) {
	userID := utils.GenerateUUID()
	productID := utils.GenerateUUID()
	thumbnailID := utils.GenerateUUID()

	type args struct {
		payload *model.UpdateProductPayload
	}
	type mockSelect struct {
		product *model.Product
		err     error
	}
	type mockGetObjectByID struct {
		res *storagePB.Object
		err error
	}
	type mockUpdate struct {
		product *model.Product
		err     error
	}
	type mockAuth struct {
		hasAccess bool
		err       error
	}

	tests := []struct {
		name              string
		args              args
		userID            string
		mockSelect        *mockSelect
		mockGetObjectByID *mockGetObjectByID
		mockUpdate        *mockUpdate
		mockAuth          *mockAuth
		want              *model.Product
		wantErr           bool
	}{
		{
			name: "success",
			args: args{
				payload: &model.UpdateProductPayload{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
				},
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: &model.Product{
					ID:          productID,
					Name:        "product-1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
				},
				err: nil,
			},
			mockGetObjectByID: &mockGetObjectByID{
				res: &storagePB.Object{
					Id:       thumbnailID,
					FileName: "test.png",
					Type:     model.ThumbnailType,
					IsPublic: true,
				},
			},
			mockUpdate: &mockUpdate{
				product: &model.Product{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
				},
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			want: &model.Product{
				ID:          productID,
				Name:        "updated product",
				Description: "updated product",
				Price:       10.10,
				ThumbnailID: thumbnailID,
			},
			wantErr: false,
		},
		{
			name: "success update other user product with full access permission",
			args: args{
				payload: &model.UpdateProductPayload{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
				},
			},
			userID: utils.GenerateUUID(),
			mockSelect: &mockSelect{
				product: &model.Product{
					ID:          productID,
					Name:        "product-1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
					OwnerID:     userID,
				},
				err: nil,
			},
			mockGetObjectByID: &mockGetObjectByID{
				res: &storagePB.Object{
					Id:       thumbnailID,
					FileName: "test.png",
					Type:     model.ThumbnailType,
					IsPublic: true,
				},
			},
			mockUpdate: &mockUpdate{
				product: &model.Product{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
					OwnerID:     userID,
				},
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			want: &model.Product{
				ID:          productID,
				Name:        "updated product",
				Description: "updated product",
				Price:       10.10,
				ThumbnailID: thumbnailID,
				OwnerID:     userID,
			},
			wantErr: false,
		},
		{
			name: "error update other user product",
			args: args{
				payload: &model.UpdateProductPayload{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
				},
			},
			userID: utils.GenerateUUID(),
			mockSelect: &mockSelect{
				product: &model.Product{
					ID:          productID,
					Name:        "product-1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
					OwnerID:     utils.GenerateUUID(),
				},
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: false,
				err:       nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "db error when update",
			args: args{
				payload: &model.UpdateProductPayload{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
				},
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: &model.Product{
					ID:          productID,
					Name:        "product-1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
				},
				err: nil,
			},
			mockGetObjectByID: &mockGetObjectByID{
				res: &storagePB.Object{
					Id:       thumbnailID,
					FileName: "test.png",
					Type:     model.ThumbnailType,
					IsPublic: true,
				},
			},
			mockUpdate: &mockUpdate{
				product: &model.Product{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
				},
				err: errors.New("db error"),
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "product not found",
			args: args{
				payload: &model.UpdateProductPayload{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
				},
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: nil,
				err:     nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "error when find product",
			args: args{
				payload: &model.UpdateProductPayload{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
				},
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: nil,
				err:     errors.New("db error"),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "permission denied",
			args: args{
				payload: &model.UpdateProductPayload{
					ID:          productID,
					Name:        "updated product",
					Description: "updated product",
					Price:       10.10,
					ThumbnailID: thumbnailID,
				},
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: &model.Product{
					ID:          productID,
					Name:        "product-1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
					OwnerID:     utils.GenerateUUID(),
				},
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: false,
				err:       nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			ctx = context.WithValue(ctx, constant.KeyUserIDCtx, tt.userID)

			uc := NewProductUsecase()
			db, _ := utils.NewDBMock()
			err := uc.InjectDB(db)
			utils.ContinueOrFatal(err)
			mockAuthClient := authMock.NewMockAuthServiceClient(ctrl)
			err = uc.InjectAuthClient(mockAuthClient)
			utils.ContinueOrFatal(err)
			mockStorageClient := storageMock.NewMockStorageServiceClient(ctrl)
			err = uc.InjectStorageClient(mockStorageClient)
			utils.ContinueOrFatal(err)
			mockProductRepo := mock.NewMockProductRepository(ctrl)
			err = uc.InjectProductRepo(mockProductRepo)
			utils.ContinueOrFatal(err)

			if tt.mockAuth != nil {
				mockAuthClient.EXPECT().
					HasAccess(gomock.Any(), gomock.Any()).
					Times(1).Return(&wrapperspb.BoolValue{
					Value: tt.mockAuth.hasAccess,
				}, tt.mockAuth.err)
			}

			if tt.mockSelect != nil {
				mockProductRepo.EXPECT().FindByID(gomock.Any(), tt.args.payload.ID).Times(1).Return(tt.mockSelect.product, tt.mockSelect.err)
			}

			if tt.mockGetObjectByID != nil {
				mockStorageClient.EXPECT().GetObjectByID(gomock.Any(), &storagePB.GetObjectByIDRequest{
					UserId:   tt.userID,
					ObjectId: tt.args.payload.ThumbnailID,
				}).Times(1).Return(tt.mockGetObjectByID.res, tt.mockGetObjectByID.err)
			}

			if tt.mockUpdate != nil {
				mockProductRepo.EXPECT().Update(gomock.Any(), tt.mockUpdate.product).Times(1).Return(tt.mockUpdate.err)
			}

			got, err := uc.Update(ctx, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("productUsecase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productUsecase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_productUsecase_Delete(t *testing.T) {
	userID := utils.GenerateUUID()
	productID := utils.GenerateUUID()
	thumbnailID := utils.GenerateUUID()

	type args struct {
		id string
	}
	type mockSelect struct {
		product *model.Product
		err     error
	}
	type mockDelete struct {
		err error
	}
	type mockAuth struct {
		hasAccess bool
		err       error
	}

	tests := []struct {
		name       string
		args       args
		userID     string
		mockSelect *mockSelect
		mockDelete *mockDelete
		mockAuth   *mockAuth
		wantErr    bool
	}{
		{
			name: "success",
			args: args{
				id: productID,
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: &model.Product{
					ID:          productID,
					Name:        "product-1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
					OwnerID:     userID,
				},
				err: nil,
			},
			mockDelete: &mockDelete{
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			wantErr: false,
		},
		{
			name: "error when delete product",
			args: args{
				id: productID,
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: &model.Product{
					ID:          productID,
					Name:        "product-1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
					OwnerID:     userID,
				},
				err: nil,
			},
			mockDelete: &mockDelete{
				err: errors.New("db error"),
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			wantErr: true,
		},
		{
			name: "permission denied",
			args: args{
				id: productID,
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: &model.Product{
					ID:          productID,
					Name:        "product-1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
					OwnerID:     utils.GenerateUUID(),
				},
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: false,
				err:       nil,
			},
			wantErr: true,
		},
		{
			name: "error when find product",
			args: args{
				id: productID,
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: nil,
				err:     errors.New("db error"),
			},
			wantErr: true,
		},
		{
			name: "error when product not found",
			args: args{
				id: productID,
			},
			userID: userID,
			mockSelect: &mockSelect{
				product: nil,
				err:     nil,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			ctx = context.WithValue(ctx, constant.KeyUserIDCtx, tt.userID)

			uc := NewProductUsecase()
			db, _ := utils.NewDBMock()
			err := uc.InjectDB(db)
			utils.ContinueOrFatal(err)
			mockAuthClient := authMock.NewMockAuthServiceClient(ctrl)
			err = uc.InjectAuthClient(mockAuthClient)
			utils.ContinueOrFatal(err)
			mockProductRepo := mock.NewMockProductRepository(ctrl)
			err = uc.InjectProductRepo(mockProductRepo)
			utils.ContinueOrFatal(err)

			if tt.mockSelect != nil {
				mockProductRepo.EXPECT().FindByID(gomock.Any(), tt.args.id).Times(1).Return(tt.mockSelect.product, tt.mockSelect.err)
				if tt.mockAuth != nil && tt.mockSelect.product.OwnerID != tt.userID {
					mockAuthClient.EXPECT().HasAccess(gomock.Any(), gomock.Any()).Times(1).Return(&wrapperspb.BoolValue{
						Value: tt.mockAuth.hasAccess,
					}, tt.mockAuth.err)
				}
			}

			if tt.mockDelete != nil {
				mockProductRepo.EXPECT().DeleteByID(gomock.Any(), tt.args.id).Times(1).Return(tt.mockDelete.err)
			}

			if err := uc.Delete(ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("productUsecase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_productUsecase_FindPaginatedIDs(t *testing.T) {
	userID := utils.GenerateUUID()
	productID := utils.GenerateUUID()

	type args struct {
		datasource int
		req        *model.PaginationPayload
	}
	type mockFindPaginatedIDs struct {
		ids   []string
		count int64
		err   error
	}
	type mockAuth struct {
		hasAccess bool
		err       error
	}
	tests := []struct {
		name                 string
		args                 args
		userID               string
		mockFindPaginatedIDs *mockFindPaginatedIDs
		mockAuth             *mockAuth
		want                 *model.PaginationResponse
		wantErr              bool
	}{
		{
			name: "success",
			args: args{
				datasource: constant.SourceDB,
				req: &model.PaginationPayload{
					Search: "",
					Sort:   []string{},
					Limit:  10,
					Page:   1,
				},
			},
			userID: userID,
			mockFindPaginatedIDs: &mockFindPaginatedIDs{
				ids:   []string{productID},
				count: 1,
				err:   nil,
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			want: &model.PaginationResponse{
				Meta: &model.PaginationPayload{
					Search: "",
					Sort:   []string{},
					Limit:  10,
					Page:   1,
				},
				Count:   1,
				MaxPage: 1,
				Items:   []string{productID},
			},
			wantErr: false,
		},
		{
			name: "error when get data",
			args: args{
				datasource: constant.SourceDB,
				req: &model.PaginationPayload{
					Search: "",
					Sort:   []string{},
					Limit:  10,
					Page:   1,
				},
			},
			userID: userID,
			mockFindPaginatedIDs: &mockFindPaginatedIDs{
				ids:   []string{},
				count: 0,
				err:   errors.New("db error"),
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "permission denied",
			args: args{
				datasource: constant.SourceDB,
				req: &model.PaginationPayload{
					Search: "",
					Sort:   []string{},
					Limit:  10,
					Page:   1,
				},
			},
			userID: userID,
			mockAuth: &mockAuth{
				hasAccess: false,
				err:       nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			ctx = context.WithValue(ctx, constant.KeyUserIDCtx, tt.userID)
			ctx = context.WithValue(ctx, constant.KeyDataSource, tt.args.datasource)

			uc := NewProductUsecase()
			db, _ := utils.NewDBMock()
			err := uc.InjectDB(db)
			utils.ContinueOrFatal(err)
			mockAuthClient := authMock.NewMockAuthServiceClient(ctrl)
			err = uc.InjectAuthClient(mockAuthClient)
			utils.ContinueOrFatal(err)
			mockProductRepo := mock.NewMockProductRepository(ctrl)
			err = uc.InjectProductRepo(mockProductRepo)
			utils.ContinueOrFatal(err)

			if tt.mockAuth != nil {
				mockAuthClient.EXPECT().HasAccess(gomock.Any(), gomock.Any()).Times(1).Return(&wrapperspb.BoolValue{
					Value: tt.mockAuth.hasAccess,
				}, tt.mockAuth.err)
			}

			if tt.mockFindPaginatedIDs != nil {
				mockProductRepo.EXPECT().FindPaginatedIDs(gomock.Any(), tt.args.req).Times(1).Return(tt.mockFindPaginatedIDs.ids, tt.mockFindPaginatedIDs.count, tt.mockFindPaginatedIDs.err)
			}

			got, err := uc.FindPaginatedIDs(ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("productUsecase.FindPaginatedIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productUsecase.FindPaginatedIDs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_productUsecase_FindByID(t *testing.T) {
	userID := utils.GenerateUUID()
	productID := utils.GenerateUUID()
	thumbnailID := utils.GenerateUUID()
	type args struct {
		id string
	}
	type mockFindByID struct {
		product *model.Product
		err     error
	}
	type mockAuth struct {
		hasAccess bool
		err       error
	}
	tests := []struct {
		name         string
		args         args
		userID       string
		mockFindByID *mockFindByID
		mockAuth     *mockAuth
		want         *model.Product
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				id: productID,
			},
			userID: userID,
			mockFindByID: &mockFindByID{
				product: &model.Product{
					ID:          productID,
					Name:        "product1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
					OwnerID:     userID,
				},
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			want: &model.Product{
				ID:          productID,
				Name:        "product1",
				Description: "product1",
				Price:       17.17,
				ThumbnailID: thumbnailID,
				OwnerID:     userID,
			},
			wantErr: false,
		},
		{
			name: "error when get data",
			args: args{
				id: productID,
			},
			userID: userID,
			mockFindByID: &mockFindByID{
				product: nil,
				err:     errors.New("db error"),
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "permission denied",
			args: args{
				id: productID,
			},
			userID: userID,
			mockAuth: &mockAuth{
				hasAccess: false,
				err:       nil,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			ctx = context.WithValue(ctx, constant.KeyUserIDCtx, tt.userID)

			uc := NewProductUsecase()
			db, _ := utils.NewDBMock()
			err := uc.InjectDB(db)
			utils.ContinueOrFatal(err)
			mockAuthClient := authMock.NewMockAuthServiceClient(ctrl)
			err = uc.InjectAuthClient(mockAuthClient)
			utils.ContinueOrFatal(err)
			mockProductRepo := mock.NewMockProductRepository(ctrl)
			err = uc.InjectProductRepo(mockProductRepo)
			utils.ContinueOrFatal(err)

			if tt.mockAuth != nil {
				mockAuthClient.EXPECT().HasAccess(gomock.Any(), gomock.Any()).Times(1).Return(&wrapperspb.BoolValue{
					Value: tt.mockAuth.hasAccess,
				}, tt.mockAuth.err)
			}

			if tt.mockFindByID != nil {
				mockProductRepo.EXPECT().FindByID(gomock.Any(), tt.args.id).Times(1).Return(tt.mockFindByID.product, tt.mockFindByID.err)
			}

			got, err := uc.FindByID(ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("productUsecase.FindByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productUsecase.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_productUsecase_FindByIDs(t *testing.T) {
	userID := utils.GenerateUUID()
	productID := utils.GenerateUUID()
	thumbnailID := utils.GenerateUUID()
	type args struct {
		ids []string
	}
	type mockFindByID struct {
		product *model.Product
		err     error
	}
	type mockAuth struct {
		hasAccess bool
		err       error
	}
	tests := []struct {
		name         string
		args         args
		userID       string
		mockFindByID *mockFindByID
		mockAuth     *mockAuth
		want         model.Products
		wantErr      bool
	}{
		{
			name: "success",
			args: args{
				ids: []string{productID},
			},
			userID: userID,
			mockFindByID: &mockFindByID{
				product: &model.Product{
					ID:          productID,
					Name:        "product1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
					OwnerID:     userID,
				},
				err: nil,
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			want: model.Products{
				{
					ID:          productID,
					Name:        "product1",
					Description: "product1",
					Price:       17.17,
					ThumbnailID: thumbnailID,
					OwnerID:     userID,
				},
			},
			wantErr: false,
		},
		{
			name: "ignore error if data not found",
			args: args{
				ids: []string{productID},
			},
			userID: userID,
			mockFindByID: &mockFindByID{
				product: nil,
				err:     errors.New("db error"),
			},
			mockAuth: &mockAuth{
				hasAccess: true,
				err:       nil,
			},
			want:    model.Products{},
			wantErr: false,
		},
		{
			name: "permission denied",
			args: args{
				ids: []string{productID},
			},
			userID: userID,
			mockAuth: &mockAuth{
				hasAccess: false,
				err:       nil,
			},
			want:    model.Products{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.TODO()
			ctx = context.WithValue(ctx, constant.KeyUserIDCtx, tt.userID)

			uc := NewProductUsecase()
			db, _ := utils.NewDBMock()
			err := uc.InjectDB(db)
			utils.ContinueOrFatal(err)
			mockAuthClient := authMock.NewMockAuthServiceClient(ctrl)
			err = uc.InjectAuthClient(mockAuthClient)
			utils.ContinueOrFatal(err)
			mockProductRepo := mock.NewMockProductRepository(ctrl)
			err = uc.InjectProductRepo(mockProductRepo)
			utils.ContinueOrFatal(err)

			wg := sync.WaitGroup{}
			if tt.mockAuth != nil {
				mockAuthClient.EXPECT().HasAccess(gomock.Any(), gomock.Any()).Times(1).Return(&wrapperspb.BoolValue{
					Value: tt.mockAuth.hasAccess,
				}, tt.mockAuth.err)
			}

			if tt.mockFindByID != nil {
				wg.Add(1)
				mockProductRepo.EXPECT().FindByID(gomock.Any(), tt.args.ids[0]).Times(1).DoAndReturn(func(_ context.Context, id string) (*model.Product, error) {
					defer wg.Done()
					return tt.mockFindByID.product, tt.mockFindByID.err
				})
			}

			got, err := uc.FindByIDs(ctx, tt.args.ids)
			if (err != nil) != tt.wantErr {
				t.Errorf("productUsecase.FindByIDs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("productUsecase.FindByIDs() = %v, want %v", got, tt.want)
			}
			wg.Wait()
		})
	}
}
