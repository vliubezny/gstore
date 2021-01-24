package service

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/storage"
)

var (
	ctx     = context.Background()
	errTest = errors.New("test")
)

func TestService_GetCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	categories := []model.Category{
		{ID: 1, Name: "Test1"},
		{ID: 2, Name: "Test2"},
	}

	st := storage.NewMockStorage(ctrl)
	st.EXPECT().GetCategories(ctx).Return(categories, nil)

	s := New(st)

	cs, err := s.GetCategories(ctx)
	assert.NoError(t, err)
	assert.Equal(t, categories, cs)
}

func TestService_GetCategories_Err(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	st := storage.NewMockStorage(ctrl)
	st.EXPECT().GetCategories(ctx).Return(nil, errTest)

	s := New(st)

	cs, err := s.GetCategories(ctx)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, errTest))
	assert.Nil(t, cs)
}

func TestService_GetCategory(t *testing.T) {
	testCases := []struct {
		desc      string
		rCategory model.Category
		rErr      error
		category  model.Category
		err       error
	}{
		{
			desc:      "success",
			rCategory: model.Category{ID: 1, Name: "Test1"},
			rErr:      nil,
			category:  model.Category{ID: 1, Name: "Test1"},
			err:       nil,
		},
		{
			desc:      "ErrNotFound",
			rCategory: model.Category{},
			rErr:      storage.ErrNotFound,
			category:  model.Category{},
			err:       ErrNotFound,
		},
		{
			desc:      "unexpected error",
			rCategory: model.Category{},
			rErr:      errTest,
			category:  model.Category{},
			err:       errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := int64(1)

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().GetCategory(ctx, id).Return(tC.rCategory, tC.rErr)

			s := New(st)

			category, err := s.GetCategory(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.category, category)
		})
	}
}

func TestService_CreateCategory(t *testing.T) {
	testCases := []struct {
		desc      string
		rCategory model.Category
		rErr      error
		category  model.Category
		err       error
	}{
		{
			desc:      "success",
			rCategory: model.Category{ID: 1, Name: "Test1"},
			rErr:      nil,
			category:  model.Category{Name: "Test1"},
			err:       nil,
		},
		{
			desc:      "unexpected error",
			rCategory: model.Category{},
			rErr:      errTest,
			category:  model.Category{ID: 1, Name: "Test1"},
			err:       errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().CreateCategory(ctx, tC.category).Return(tC.rCategory, tC.rErr)

			s := New(st)

			c, err := s.CreateCategory(ctx, tC.category)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.rCategory, c)
		})
	}
}

func TestService_UpdateCategory(t *testing.T) {
	testCases := []struct {
		desc     string
		rErr     error
		category model.Category
		err      error
	}{
		{
			desc:     "success",
			rErr:     nil,
			category: model.Category{ID: 1, Name: "Test1"},
			err:      nil,
		},
		{
			desc:     "ErrNotFound",
			rErr:     storage.ErrNotFound,
			category: model.Category{ID: 1, Name: "Test1"},
			err:      ErrNotFound,
		},
		{
			desc:     "unexpected error",
			rErr:     errTest,
			category: model.Category{ID: 1, Name: "Test1"},
			err:      errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().UpdateCategory(ctx, tC.category).Return(tC.rErr)

			s := New(st)

			err := s.UpdateCategory(ctx, tC.category)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_DeleteCategory(t *testing.T) {
	testCases := []struct {
		desc string
		rErr error
		err  error
	}{
		{
			desc: "success",
			rErr: nil,
			err:  nil,
		},
		{
			desc: "ErrNotFound",
			rErr: storage.ErrNotFound,
			err:  ErrNotFound,
		},
		{
			desc: "unexpected error",
			rErr: errTest,
			err:  errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := int64(1)

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().DeleteCategory(ctx, id).Return(tC.rErr)

			s := New(st)

			err := s.DeleteCategory(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_GetStores(t *testing.T) {
	testStores := []*model.Store{
		{ID: 1, Name: "AAA"},
		{ID: 2, Name: "BBB"},
	}

	testCases := []struct {
		desc    string
		rStores []*model.Store
		rErr    error
		stores  []*model.Store
		err     error
	}{
		{
			desc:    "success",
			rStores: testStores,
			rErr:    nil,
			stores:  testStores,
			err:     nil,
		},
		{
			desc:    "unexpected error",
			rStores: nil,
			rErr:    errTest,
			stores:  nil,
			err:     errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().GetStores(ctx).Return(tC.rStores, tC.rErr)

			s := New(st)

			stores, err := s.GetStores(ctx)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.stores, stores)
		})
	}
}

func TestService_GetStore(t *testing.T) {
	testCases := []struct {
		desc   string
		rStore *model.Store
		rErr   error
		store  *model.Store
		err    error
	}{
		{
			desc:   "success",
			rStore: &model.Store{ID: 1, Name: "Test1"},
			rErr:   nil,
			store:  &model.Store{ID: 1, Name: "Test1"},
			err:    nil,
		},
		{
			desc:   "ErrNotFound",
			rStore: nil,
			rErr:   storage.ErrNotFound,
			store:  nil,
			err:    ErrNotFound,
		},
		{
			desc:   "unexpected error",
			rStore: nil,
			rErr:   errTest,
			store:  nil,
			err:    errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := int64(1)

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().GetStore(ctx, id).Return(tC.rStore, tC.rErr)

			s := New(st)

			store, err := s.GetStore(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.store, store)
		})
	}
}

func TestService_CreateStore(t *testing.T) {
	testCases := []struct {
		desc  string
		rErr  error
		store *model.Store
		err   error
	}{
		{
			desc:  "success",
			rErr:  nil,
			store: &model.Store{ID: 1, Name: "Test1"},
			err:   nil,
		},
		{
			desc:  "unexpected error",
			rErr:  errTest,
			store: &model.Store{ID: 1, Name: "Test1"},
			err:   errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().CreateStore(ctx, tC.store).Return(tC.rErr)

			s := New(st)

			err := s.CreateStore(ctx, tC.store)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_UpdateStore(t *testing.T) {
	testCases := []struct {
		desc  string
		rErr  error
		store *model.Store
		err   error
	}{
		{
			desc:  "success",
			rErr:  nil,
			store: &model.Store{ID: 1, Name: "Test1"},
			err:   nil,
		},
		{
			desc:  "ErrNotFound",
			rErr:  storage.ErrNotFound,
			store: &model.Store{ID: 1, Name: "Test1"},
			err:   ErrNotFound,
		},
		{
			desc:  "unexpected error",
			rErr:  errTest,
			store: &model.Store{ID: 1, Name: "Test1"},
			err:   errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().UpdateStore(ctx, tC.store).Return(tC.rErr)

			s := New(st)

			err := s.UpdateStore(ctx, tC.store)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_DeleteStore(t *testing.T) {
	testCases := []struct {
		desc string
		rErr error
		err  error
	}{
		{
			desc: "success",
			rErr: nil,
			err:  nil,
		},
		{
			desc: "ErrNotFound",
			rErr: storage.ErrNotFound,
			err:  ErrNotFound,
		},
		{
			desc: "unexpected error",
			rErr: errTest,
			err:  errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := int64(1)

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().DeleteStore(ctx, id).Return(tC.rErr)

			s := New(st)

			err := s.DeleteStore(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_GetProducts(t *testing.T) {
	testProducts := []*model.Product{
		{ID: 1, CategoryID: 1, Name: "AAA", Description: "D-AAA"},
		{ID: 2, CategoryID: 1, Name: "BBB", Description: "D-BBB"},
	}

	testCases := []struct {
		desc      string
		rProducts []*model.Product
		rErr      error
		products  []*model.Product
		err       error
	}{
		{
			desc:      "success",
			rProducts: testProducts,
			rErr:      nil,
			products:  testProducts,
			err:       nil,
		},
		{
			desc:      "unexpected error",
			rProducts: nil,
			rErr:      errTest,
			products:  nil,
			err:       errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().GetProducts(ctx, int64(1)).Return(tC.rProducts, tC.rErr)

			s := New(st)

			stores, err := s.GetProducts(ctx, 1)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.products, stores)
		})
	}
}

func TestService_GetProduct(t *testing.T) {
	testCases := []struct {
		desc     string
		rProduct *model.Product
		rErr     error
		product  *model.Product
		err      error
	}{
		{
			desc:     "success",
			rProduct: &model.Product{ID: 1, CategoryID: 1, Name: "Test1", Description: "1 test"},
			rErr:     nil,
			product:  &model.Product{ID: 1, CategoryID: 1, Name: "Test1", Description: "1 test"},
			err:      nil,
		},
		{
			desc:     "ErrNotFound",
			rProduct: nil,
			rErr:     storage.ErrNotFound,
			product:  nil,
			err:      ErrNotFound,
		},
		{
			desc:     "unexpected error",
			rProduct: nil,
			rErr:     errTest,
			product:  nil,
			err:      errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := int64(1)

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().GetProduct(ctx, id).Return(tC.rProduct, tC.rErr)

			s := New(st)

			product, err := s.GetProduct(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.product, product)
		})
	}
}

func TestService_CreateProduct(t *testing.T) {
	testCases := []struct {
		desc    string
		rErr    error
		product *model.Product
		err     error
	}{
		{
			desc:    "success",
			rErr:    nil,
			product: &model.Product{ID: 1, CategoryID: 1, Name: "Test1", Description: "1 test"},
			err:     nil,
		},
		{
			desc:    "unexpected error",
			rErr:    errTest,
			product: &model.Product{ID: 1, CategoryID: 1, Name: "Test1", Description: "1 test"},
			err:     errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().CreateProduct(ctx, tC.product).Return(tC.rErr)

			s := New(st)

			err := s.CreateProduct(ctx, tC.product)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_UpdateProduct(t *testing.T) {
	testCases := []struct {
		desc    string
		rErr    error
		product *model.Product
		err     error
	}{
		{
			desc:    "success",
			rErr:    nil,
			product: &model.Product{ID: 1, CategoryID: 1, Name: "Test1", Description: "1 test"},
			err:     nil,
		},
		{
			desc:    "ErrNotFound",
			rErr:    storage.ErrNotFound,
			product: &model.Product{ID: 1, CategoryID: 1, Name: "Test1", Description: "1 test"},
			err:     ErrNotFound,
		},
		{
			desc:    "unexpected error",
			rErr:    errTest,
			product: &model.Product{ID: 1, CategoryID: 1, Name: "Test1", Description: "1 test"},
			err:     errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().UpdateProduct(ctx, tC.product).Return(tC.rErr)

			s := New(st)

			err := s.UpdateProduct(ctx, tC.product)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_DeleteProduct(t *testing.T) {
	testCases := []struct {
		desc string
		rErr error
		err  error
	}{
		{
			desc: "success",
			rErr: nil,
			err:  nil,
		},
		{
			desc: "ErrNotFound",
			rErr: storage.ErrNotFound,
			err:  ErrNotFound,
		},
		{
			desc: "unexpected error",
			rErr: errTest,
			err:  errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := int64(1)

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().DeleteProduct(ctx, id).Return(tC.rErr)

			s := New(st)

			err := s.DeleteProduct(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}
