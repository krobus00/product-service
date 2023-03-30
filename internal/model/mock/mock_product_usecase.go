// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/krobus00/product-service/internal/model (interfaces: ProductUsecase)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	auth "github.com/krobus00/auth-service/pb/auth"
	model "github.com/krobus00/product-service/internal/model"
	storage "github.com/krobus00/storage-service/pb/storage"
	gorm "gorm.io/gorm"
)

// MockProductUsecase is a mock of ProductUsecase interface.
type MockProductUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockProductUsecaseMockRecorder
}

// MockProductUsecaseMockRecorder is the mock recorder for MockProductUsecase.
type MockProductUsecaseMockRecorder struct {
	mock *MockProductUsecase
}

// NewMockProductUsecase creates a new mock instance.
func NewMockProductUsecase(ctrl *gomock.Controller) *MockProductUsecase {
	mock := &MockProductUsecase{ctrl: ctrl}
	mock.recorder = &MockProductUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProductUsecase) EXPECT() *MockProductUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockProductUsecase) Create(arg0 context.Context, arg1 *model.CreateProductPayload) (*model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockProductUsecaseMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockProductUsecase)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockProductUsecase) Delete(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockProductUsecaseMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockProductUsecase)(nil).Delete), arg0, arg1)
}

// FindByID mocks base method.
func (m *MockProductUsecase) FindByID(arg0 context.Context, arg1 string) (*model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", arg0, arg1)
	ret0, _ := ret[0].(*model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockProductUsecaseMockRecorder) FindByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockProductUsecase)(nil).FindByID), arg0, arg1)
}

// FindByIDs mocks base method.
func (m *MockProductUsecase) FindByIDs(arg0 context.Context, arg1 []string) (model.Products, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByIDs", arg0, arg1)
	ret0, _ := ret[0].(model.Products)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByIDs indicates an expected call of FindByIDs.
func (mr *MockProductUsecaseMockRecorder) FindByIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByIDs", reflect.TypeOf((*MockProductUsecase)(nil).FindByIDs), arg0, arg1)
}

// FindPaginatedIDs mocks base method.
func (m *MockProductUsecase) FindPaginatedIDs(arg0 context.Context, arg1 *model.PaginationPayload) (*model.PaginationResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindPaginatedIDs", arg0, arg1)
	ret0, _ := ret[0].(*model.PaginationResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindPaginatedIDs indicates an expected call of FindPaginatedIDs.
func (mr *MockProductUsecaseMockRecorder) FindPaginatedIDs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindPaginatedIDs", reflect.TypeOf((*MockProductUsecase)(nil).FindPaginatedIDs), arg0, arg1)
}

// InjectAuthClient mocks base method.
func (m *MockProductUsecase) InjectAuthClient(arg0 auth.AuthServiceClient) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectAuthClient", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectAuthClient indicates an expected call of InjectAuthClient.
func (mr *MockProductUsecaseMockRecorder) InjectAuthClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectAuthClient", reflect.TypeOf((*MockProductUsecase)(nil).InjectAuthClient), arg0)
}

// InjectDB mocks base method.
func (m *MockProductUsecase) InjectDB(arg0 *gorm.DB) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectDB", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectDB indicates an expected call of InjectDB.
func (mr *MockProductUsecaseMockRecorder) InjectDB(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectDB", reflect.TypeOf((*MockProductUsecase)(nil).InjectDB), arg0)
}

// InjectProductRepo mocks base method.
func (m *MockProductUsecase) InjectProductRepo(arg0 model.ProductRepository) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectProductRepo", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectProductRepo indicates an expected call of InjectProductRepo.
func (mr *MockProductUsecaseMockRecorder) InjectProductRepo(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectProductRepo", reflect.TypeOf((*MockProductUsecase)(nil).InjectProductRepo), arg0)
}

// InjectStorageClient mocks base method.
func (m *MockProductUsecase) InjectStorageClient(arg0 storage.StorageServiceClient) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InjectStorageClient", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// InjectStorageClient indicates an expected call of InjectStorageClient.
func (mr *MockProductUsecaseMockRecorder) InjectStorageClient(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InjectStorageClient", reflect.TypeOf((*MockProductUsecase)(nil).InjectStorageClient), arg0)
}

// Update mocks base method.
func (m *MockProductUsecase) Update(arg0 context.Context, arg1 *model.UpdateProductPayload) (*model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockProductUsecaseMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockProductUsecase)(nil).Update), arg0, arg1)
}
