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
	categories := []*model.Category{
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
			desc:   "ErrNotFound",
			rItems: nil,
			rErr:   storage.ErrNotFound,
			items:  nil,
			err:    ErrNotFound,
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
