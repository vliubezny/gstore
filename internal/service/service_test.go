package service

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vliubezny/gstore/internal/storage"
)

func TestService_GetCategories(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	categories := []*storage.Category{
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

	ctx := context.Background()
	testErr := errors.New("test")

	st := storage.NewMockStorage(ctrl)
	st.EXPECT().GetCategories(ctx).Return(nil, testErr)

	s := New(st)

	cs, err := s.GetCategories(ctx)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, testErr))
	assert.Nil(t, cs)
}
