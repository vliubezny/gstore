package server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/vliubezny/gstore/internal/model"
	"github.com/vliubezny/gstore/internal/service"
)

var (
	errTest = errors.New("test")
	errSkip = errors.New("skip")
)

func Test_getCategoriesHandler(t *testing.T) {
	testCases := []struct {
		desc       string
		categories []*model.Category
		err        error
		rcode      int
		rdata      string
	}{
		{
			desc: "success",
			categories: []*model.Category{
				{ID: 1, Name: "Test1"},
				{ID: 2, Name: "Test2"},
			},
			err:   nil,
			rcode: http.StatusOK,
			rdata: `[{"id":1, "name":"Test1"}, {"id":2, "name":"Test2"}]`,
		},
		{
			desc:       "internal error",
			categories: nil,
			err:        errTest,
			rcode:      http.StatusInternalServerError,
			rdata:      `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			svc.EXPECT().GetCategories(gomock.Any()).Return(tC.categories, tC.err)

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodGet, "/v1/categories", "")

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_getCategoryHandler(t *testing.T) {
	testCases := []struct {
		desc     string
		category *model.Category
		id       string
		err      error
		rcode    int
		rdata    string
	}{
		{
			desc:     "success",
			category: &model.Category{ID: 1, Name: "Test1"},
			id:       "1",
			err:      nil,
			rcode:    http.StatusOK,
			rdata:    `{"id":1, "name":"Test1"}`,
		},
		{
			desc:     "invalid category ID",
			id:       "test",
			category: nil,
			err:      errSkip,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error":"invalid category ID"}`,
		},
		{
			desc:     "not found",
			id:       "1",
			category: nil,
			err:      service.ErrNotFound,
			rcode:    http.StatusNotFound,
			rdata:    `{"error":"category not found"}`,
		},
		{
			desc:     "internal error",
			category: nil,
			id:       "1",
			err:      errTest,
			rcode:    http.StatusInternalServerError,
			rdata:    `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().GetCategory(gomock.Any(), int64(1)).Return(tC.category, tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodGet, fmt.Sprintf("/v1/categories/%s", tC.id), "")

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_createCategoryHandler(t *testing.T) {
	testCases := []struct {
		desc     string
		category *model.Category
		err      error
		input    string
		rcode    int
		rdata    string
	}{
		{
			desc:     "success",
			category: &model.Category{Name: "Test1"},
			err:      nil,
			input:    `{"name": "Test1"}`,
			rcode:    http.StatusOK,
			rdata:    `{"id":1, "name":"Test1"}`,
		},
		{
			desc:     "invalid: missing name",
			category: nil,
			input:    `{}`,
			err:      errSkip,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error":"name: cannot be blank."}`,
		},
		{
			desc:     "internal error",
			category: &model.Category{Name: "Test1"},
			err:      errTest,
			input:    `{"name": "Test1"}`,
			rcode:    http.StatusInternalServerError,
			rdata:    `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().CreateCategory(gomock.Any(), tC.category).DoAndReturn(func(_ context.Context, c *model.Category) error {
					c.ID = 1
					return tC.err
				})
			}

			router := setupTestRouter(svc)
			rec, r := newTestParametersWithAuth(http.MethodPost, "/v1/categories", tC.input)

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_updateCategoryHandler(t *testing.T) {
	testCases := []struct {
		desc     string
		category *model.Category
		err      error
		id       string
		input    string
		rcode    int
		rdata    string
	}{
		{
			desc:     "success",
			category: &model.Category{ID: 1, Name: "Test1"},
			err:      nil,
			id:       "1",
			input:    `{"name": "Test1"}`,
			rcode:    http.StatusOK,
			rdata:    `{"id":1, "name":"Test1"}`,
		},
		{
			desc:     "invalid: missing name",
			category: nil,
			id:       "1",
			input:    `{}`,
			err:      errSkip,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error":"name: cannot be blank."}`,
		},
		{
			desc:     "invalid category ID",
			category: nil,
			id:       "test",
			input:    `{"name": "Test1"}`,
			err:      errSkip,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error":"invalid category ID"}`,
		},
		{
			desc:     "not found",
			category: &model.Category{ID: 1, Name: "Test1"},
			id:       "1",
			input:    `{"name": "Test1"}`,
			err:      service.ErrNotFound,
			rcode:    http.StatusNotFound,
			rdata:    `{"error":"category not found"}`,
		},
		{
			desc:     "internal error",
			category: &model.Category{ID: 1, Name: "Test1"},
			err:      errTest,
			id:       "1",
			input:    `{"name": "Test1"}`,
			rcode:    http.StatusInternalServerError,
			rdata:    `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().UpdateCategory(gomock.Any(), tC.category).Return(tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParametersWithAuth(http.MethodPut, fmt.Sprintf("/v1/categories/%s", tC.id), tC.input)

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_deleteCategoryHandler(t *testing.T) {
	testCases := []struct {
		desc  string
		id    string
		err   error
		rcode int
		rdata string
	}{
		{
			desc:  "success",
			id:    "1",
			err:   nil,
			rcode: http.StatusNoContent,
			rdata: "",
		},
		{
			desc:  "invalid category ID",
			id:    "test",
			err:   errSkip,
			rcode: http.StatusBadRequest,
			rdata: `{"error":"invalid category ID"}`,
		},
		{
			desc:  "not found",
			id:    "1",
			err:   service.ErrNotFound,
			rcode: http.StatusNotFound,
			rdata: `{"error":"category not found"}`,
		},
		{
			desc:  "internal error",
			id:    "1",
			err:   errTest,
			rcode: http.StatusInternalServerError,
			rdata: `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().DeleteCategory(gomock.Any(), int64(1)).Return(tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParametersWithAuth(http.MethodDelete, fmt.Sprintf("/v1/categories/%s", tC.id), "")

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			if tC.rdata == "" {
				assert.Empty(t, string(body))
			} else {
				assert.JSONEq(t, tC.rdata, string(body))
			}
		})
	}
}

func Test_getStoresHandler(t *testing.T) {
	testCases := []struct {
		desc   string
		stores []*model.Store
		err    error
		rcode  int
		rdata  string
	}{
		{
			desc: "success",
			stores: []*model.Store{
				{ID: 1, Name: "Test1"},
				{ID: 2, Name: "Test2"},
			},
			err:   nil,
			rcode: http.StatusOK,
			rdata: `[{"id":1, "name":"Test1"}, {"id":2, "name":"Test2"}]`,
		},
		{
			desc:   "internal error",
			stores: nil,
			err:    errTest,
			rcode:  http.StatusInternalServerError,
			rdata:  `{"error":"internal error"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			svc.EXPECT().GetStores(gomock.Any()).Return(tC.stores, tC.err)

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodGet, "/v1/stores", "")

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_getStoreItemsHandler(t *testing.T) {
	testCases := []struct {
		desc    string
		storeID string
		items   []*model.Item
		err     error
		rcode   int
		rdata   string
	}{
		{
			desc:    "success",
			storeID: "1",
			items: []*model.Item{
				{ID: 1, StoreID: 1, Name: "Test1", Description: "Desc 1", Price: 1000},
				{ID: 2, StoreID: 1, Name: "Test2", Description: "Desc 2", Price: 2000},
			},
			err:   nil,
			rcode: http.StatusOK,
			rdata: `[{"id":1, "storeId":1, "name":"Test1", "description":"Desc 1", "price":1000},
				{"id":2, "storeId":1, "name":"Test2", "description":"Desc 2", "price":2000}]`,
		},
		{
			desc:    "internal error",
			storeID: "1",
			items:   nil,
			err:     errTest,
			rcode:   http.StatusInternalServerError,
			rdata:   `{"error":"internal error"}`,
		},
		{
			desc:    "empty store ID",
			storeID: "",
			items:   []*model.Item{},
			err:     errSkip,
			rcode:   http.StatusBadRequest,
			rdata:   `{"error":"invalid store ID"}`,
		},
		{
			desc:    "invalid store ID",
			storeID: "test",
			items:   []*model.Item{},
			err:     errSkip,
			rcode:   http.StatusBadRequest,
			rdata:   `{"error":"invalid store ID"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().GetStoreItems(gomock.Any(), int64(1)).Return(tC.items, tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodGet, fmt.Sprintf("/v1/stores/%s/items", tC.storeID), "")

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}
