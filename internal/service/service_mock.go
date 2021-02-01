// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	model "github.com/vliubezny/gstore/internal/model"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetCategories mocks base method
func (m *MockService) GetCategories(ctx context.Context) ([]model.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategories", ctx)
	ret0, _ := ret[0].([]model.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategories indicates an expected call of GetCategories
func (mr *MockServiceMockRecorder) GetCategories(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategories", reflect.TypeOf((*MockService)(nil).GetCategories), ctx)
}

// GetCategory mocks base method
func (m *MockService) GetCategory(ctx context.Context, categoryID int64) (model.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCategory", ctx, categoryID)
	ret0, _ := ret[0].(model.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCategory indicates an expected call of GetCategory
func (mr *MockServiceMockRecorder) GetCategory(ctx, categoryID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCategory", reflect.TypeOf((*MockService)(nil).GetCategory), ctx, categoryID)
}

// CreateCategory mocks base method
func (m *MockService) CreateCategory(ctx context.Context, category model.Category) (model.Category, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCategory", ctx, category)
	ret0, _ := ret[0].(model.Category)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCategory indicates an expected call of CreateCategory
func (mr *MockServiceMockRecorder) CreateCategory(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCategory", reflect.TypeOf((*MockService)(nil).CreateCategory), ctx, category)
}

// UpdateCategory mocks base method
func (m *MockService) UpdateCategory(ctx context.Context, category model.Category) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCategory", ctx, category)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCategory indicates an expected call of UpdateCategory
func (mr *MockServiceMockRecorder) UpdateCategory(ctx, category interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCategory", reflect.TypeOf((*MockService)(nil).UpdateCategory), ctx, category)
}

// DeleteCategory mocks base method
func (m *MockService) DeleteCategory(ctx context.Context, categoryID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCategory", ctx, categoryID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCategory indicates an expected call of DeleteCategory
func (mr *MockServiceMockRecorder) DeleteCategory(ctx, categoryID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCategory", reflect.TypeOf((*MockService)(nil).DeleteCategory), ctx, categoryID)
}

// GetStores mocks base method
func (m *MockService) GetStores(ctx context.Context) ([]*model.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStores", ctx)
	ret0, _ := ret[0].([]*model.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStores indicates an expected call of GetStores
func (mr *MockServiceMockRecorder) GetStores(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStores", reflect.TypeOf((*MockService)(nil).GetStores), ctx)
}

// GetStore mocks base method
func (m *MockService) GetStore(ctx context.Context, storeID int64) (*model.Store, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStore", ctx, storeID)
	ret0, _ := ret[0].(*model.Store)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStore indicates an expected call of GetStore
func (mr *MockServiceMockRecorder) GetStore(ctx, storeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStore", reflect.TypeOf((*MockService)(nil).GetStore), ctx, storeID)
}

// CreateStore mocks base method
func (m *MockService) CreateStore(ctx context.Context, store *model.Store) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateStore", ctx, store)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateStore indicates an expected call of CreateStore
func (mr *MockServiceMockRecorder) CreateStore(ctx, store interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateStore", reflect.TypeOf((*MockService)(nil).CreateStore), ctx, store)
}

// UpdateStore mocks base method
func (m *MockService) UpdateStore(ctx context.Context, store *model.Store) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStore", ctx, store)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateStore indicates an expected call of UpdateStore
func (mr *MockServiceMockRecorder) UpdateStore(ctx, store interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStore", reflect.TypeOf((*MockService)(nil).UpdateStore), ctx, store)
}

// DeleteStore mocks base method
func (m *MockService) DeleteStore(ctx context.Context, storeID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteStore", ctx, storeID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteStore indicates an expected call of DeleteStore
func (mr *MockServiceMockRecorder) DeleteStore(ctx, storeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteStore", reflect.TypeOf((*MockService)(nil).DeleteStore), ctx, storeID)
}

// GetProducts mocks base method
func (m *MockService) GetProducts(ctx context.Context, categoryID int64) ([]*model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProducts", ctx, categoryID)
	ret0, _ := ret[0].([]*model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProducts indicates an expected call of GetProducts
func (mr *MockServiceMockRecorder) GetProducts(ctx, categoryID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProducts", reflect.TypeOf((*MockService)(nil).GetProducts), ctx, categoryID)
}

// GetProduct mocks base method
func (m *MockService) GetProduct(ctx context.Context, productID int64) (*model.Product, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProduct", ctx, productID)
	ret0, _ := ret[0].(*model.Product)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProduct indicates an expected call of GetProduct
func (mr *MockServiceMockRecorder) GetProduct(ctx, productID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProduct", reflect.TypeOf((*MockService)(nil).GetProduct), ctx, productID)
}

// CreateProduct mocks base method
func (m *MockService) CreateProduct(ctx context.Context, product *model.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateProduct", ctx, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateProduct indicates an expected call of CreateProduct
func (mr *MockServiceMockRecorder) CreateProduct(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateProduct", reflect.TypeOf((*MockService)(nil).CreateProduct), ctx, product)
}

// UpdateProduct mocks base method
func (m *MockService) UpdateProduct(ctx context.Context, product *model.Product) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProduct", ctx, product)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProduct indicates an expected call of UpdateProduct
func (mr *MockServiceMockRecorder) UpdateProduct(ctx, product interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProduct", reflect.TypeOf((*MockService)(nil).UpdateProduct), ctx, product)
}

// DeleteProduct mocks base method
func (m *MockService) DeleteProduct(ctx context.Context, productID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteProduct", ctx, productID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteProduct indicates an expected call of DeleteProduct
func (mr *MockServiceMockRecorder) DeleteProduct(ctx, productID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteProduct", reflect.TypeOf((*MockService)(nil).DeleteProduct), ctx, productID)
}

// GetStorePositions mocks base method
func (m *MockService) GetStorePositions(ctx context.Context, storeID int64) ([]model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetStorePositions", ctx, storeID)
	ret0, _ := ret[0].([]model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetStorePositions indicates an expected call of GetStorePositions
func (mr *MockServiceMockRecorder) GetStorePositions(ctx, storeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetStorePositions", reflect.TypeOf((*MockService)(nil).GetStorePositions), ctx, storeID)
}

// GetProductPositions mocks base method
func (m *MockService) GetProductPositions(ctx context.Context, productID int64) ([]model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProductPositions", ctx, productID)
	ret0, _ := ret[0].([]model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProductPositions indicates an expected call of GetProductPositions
func (mr *MockServiceMockRecorder) GetProductPositions(ctx, productID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProductPositions", reflect.TypeOf((*MockService)(nil).GetProductPositions), ctx, productID)
}

// SetPosition mocks base method
func (m *MockService) SetPosition(ctx context.Context, position model.Position) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetPosition", ctx, position)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetPosition indicates an expected call of SetPosition
func (mr *MockServiceMockRecorder) SetPosition(ctx, position interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetPosition", reflect.TypeOf((*MockService)(nil).SetPosition), ctx, position)
}

// DeletePosition mocks base method
func (m *MockService) DeletePosition(ctx context.Context, productID, storeID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeletePosition", ctx, productID, storeID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeletePosition indicates an expected call of DeletePosition
func (mr *MockServiceMockRecorder) DeletePosition(ctx, productID, storeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeletePosition", reflect.TypeOf((*MockService)(nil).DeletePosition), ctx, productID, storeID)
}
