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

func TestService_GetStoreItems(t *testing.T) {
	testItems := []*model.Item{
		{ID: 1, StoreID: 1, Name: "AAA", Description: "D-AAA", Price: 1000},
		{ID: 2, StoreID: 1, Name: "BBB", Description: "D-BBB", Price: 2000},
	}

	testCases := []struct {
		desc   string
		rItems []*model.Item
		rErr   error
		items  []*model.Item
		err    error
	}{
		{
			desc:   "success",
			rItems: testItems,
			rErr:   nil,
			items:  testItems,
			err:    nil,
		},
		{
			desc:   "unexpected error",
			rItems: nil,
			rErr:   errTest,
			items:  nil,
			err:    errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().GetStoreItems(ctx, int64(1)).Return(tC.rItems, tC.rErr)

			s := New(st)

			stores, err := s.GetStoreItems(ctx, 1)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.items, stores)
		})
	}
}

func TestService_GetStoreItem(t *testing.T) {
	testCases := []struct {
		desc  string
		rItem *model.Item
		rErr  error
		item  *model.Item
		err   error
	}{
		{
			desc:  "success",
			rItem: &model.Item{ID: 1, StoreID: 1, Name: "Test1", Description: "1 test", Price: 100},
			rErr:  nil,
			item:  &model.Item{ID: 1, StoreID: 1, Name: "Test1", Description: "1 test", Price: 100},
			err:   nil,
		},
		{
			desc:  "ErrNotFound",
			rItem: nil,
			rErr:  storage.ErrNotFound,
			item:  nil,
			err:   ErrNotFound,
		},
		{
			desc:  "unexpected error",
			rItem: nil,
			rErr:  errTest,
			item:  nil,
			err:   errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			id := int64(1)

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().GetStoreItem(ctx, id).Return(tC.rItem, tC.rErr)

			s := New(st)

			item, err := s.GetStoreItem(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
			assert.Equal(t, tC.item, item)
		})
	}
}

func TestService_CreateStoreItem(t *testing.T) {
	testCases := []struct {
		desc string
		rErr error
		item *model.Item
		err  error
	}{
		{
			desc: "success",
			rErr: nil,
			item: &model.Item{ID: 1, StoreID: 1, Name: "Test1", Description: "1 test", Price: 100},
			err:  nil,
		},
		{
			desc: "unexpected error",
			rErr: errTest,
			item: &model.Item{ID: 1, StoreID: 1, Name: "Test1", Description: "1 test", Price: 100},
			err:  errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().CreateStoreItem(ctx, tC.item).Return(tC.rErr)

			s := New(st)

			err := s.CreateStoreItem(ctx, tC.item)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_UpdateStoreItem(t *testing.T) {
	testCases := []struct {
		desc string
		rErr error
		item *model.Item
		err  error
	}{
		{
			desc: "success",
			rErr: nil,
			item: &model.Item{ID: 1, StoreID: 1, Name: "Test1", Description: "1 test", Price: 100},
			err:  nil,
		},
		{
			desc: "ErrNotFound",
			rErr: storage.ErrNotFound,
			item: &model.Item{ID: 1, StoreID: 1, Name: "Test1", Description: "1 test", Price: 100},
			err:  ErrNotFound,
		},
		{
			desc: "unexpected error",
			rErr: errTest,
			item: &model.Item{ID: 1, StoreID: 1, Name: "Test1", Description: "1 test", Price: 100},
			err:  errTest,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			st := storage.NewMockStorage(ctrl)
			st.EXPECT().UpdateStoreItem(ctx, tC.item).Return(tC.rErr)

			s := New(st)

			err := s.UpdateStoreItem(ctx, tC.item)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}

func TestService_DeleteStoreItem(t *testing.T) {
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
			st.EXPECT().DeleteStoreItem(ctx, id).Return(tC.rErr)

			s := New(st)

			err := s.DeleteStoreItem(ctx, id)
			assert.True(t, errors.Is(err, tC.err), fmt.Sprintf("wanted %s got %s", tC.err, err))
		})
	}
}
