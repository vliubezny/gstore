package server

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/shopspring/decimal"
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
		categories []model.Category
		err        error
		rcode      int
		rdata      string
	}{
		{
			desc: "success",
			categories: []model.Category{
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
		category model.Category
		id       string
		err      error
		rcode    int
		rdata    string
	}{
		{
			desc:     "success",
			category: model.Category{ID: 1, Name: "Test1"},
			id:       "1",
			err:      nil,
			rcode:    http.StatusOK,
			rdata:    `{"id":1, "name":"Test1"}`,
		},
		{
			desc:     "invalid category ID",
			id:       "test",
			category: model.Category{},
			err:      errSkip,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error":"invalid category ID"}`,
		},
		{
			desc:     "not found",
			id:       "1",
			category: model.Category{},
			err:      service.ErrNotFound,
			rcode:    http.StatusNotFound,
			rdata:    `{"error":"category not found"}`,
		},
		{
			desc:     "internal error",
			category: model.Category{},
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
		category model.Category
		err      error
		input    string
		rcode    int
		rdata    string
	}{
		{
			desc:     "success",
			category: model.Category{Name: "Test1"},
			err:      nil,
			input:    `{"name": "Test1"}`,
			rcode:    http.StatusOK,
			rdata:    `{"id":1, "name":"Test1"}`,
		},
		{
			desc:     "invalid: missing name",
			category: model.Category{},
			input:    `{}`,
			err:      errSkip,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error":"name is a required field"}`,
		},
		{
			desc:     "internal error",
			category: model.Category{Name: "Test1"},
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
				svc.EXPECT().CreateCategory(gomock.Any(), tC.category).DoAndReturn(func(_ context.Context, c model.Category) (model.Category, error) {
					c.ID = 1
					return c, tC.err
				})
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodPost, "/v1/categories", tC.input)

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
		category model.Category
		err      error
		id       string
		input    string
		rcode    int
		rdata    string
	}{
		{
			desc:     "success",
			category: model.Category{ID: 1, Name: "Test1"},
			err:      nil,
			id:       "1",
			input:    `{"name": "Test1"}`,
			rcode:    http.StatusOK,
			rdata:    `{"id":1, "name":"Test1"}`,
		},
		{
			desc:     "invalid: missing name",
			category: model.Category{},
			id:       "1",
			input:    `{}`,
			err:      errSkip,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error":"name is a required field"}`,
		},
		{
			desc:     "invalid category ID",
			category: model.Category{},
			id:       "test",
			input:    `{"name": "Test1"}`,
			err:      errSkip,
			rcode:    http.StatusBadRequest,
			rdata:    `{"error":"invalid category ID"}`,
		},
		{
			desc:     "not found",
			category: model.Category{ID: 1, Name: "Test1"},
			id:       "1",
			input:    `{"name": "Test1"}`,
			err:      service.ErrNotFound,
			rcode:    http.StatusNotFound,
			rdata:    `{"error":"category not found"}`,
		},
		{
			desc:     "internal error",
			category: model.Category{ID: 1, Name: "Test1"},
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
			rec, r := newTestParameters(http.MethodPut, fmt.Sprintf("/v1/categories/%s", tC.id), tC.input)

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
			rec, r := newTestParameters(http.MethodDelete, fmt.Sprintf("/v1/categories/%s", tC.id), "")

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
		stores []model.Store
		err    error
		rcode  int
		rdata  string
	}{
		{
			desc: "success",
			stores: []model.Store{
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

func Test_getStoreHandler(t *testing.T) {
	testCases := []struct {
		desc  string
		store model.Store
		id    string
		err   error
		rcode int
		rdata string
	}{
		{
			desc:  "success",
			store: model.Store{ID: 1, Name: "Test1"},
			id:    "1",
			err:   nil,
			rcode: http.StatusOK,
			rdata: `{"id":1, "name":"Test1"}`,
		},
		{
			desc:  "invalid store ID",
			id:    "test",
			store: model.Store{},
			err:   errSkip,
			rcode: http.StatusBadRequest,
			rdata: `{"error":"invalid store ID"}`,
		},
		{
			desc:  "not found",
			id:    "1",
			store: model.Store{},
			err:   service.ErrNotFound,
			rcode: http.StatusNotFound,
			rdata: `{"error":"store not found"}`,
		},
		{
			desc:  "internal error",
			store: model.Store{},
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
				svc.EXPECT().GetStore(gomock.Any(), int64(1)).Return(tC.store, tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodGet, fmt.Sprintf("/v1/stores/%s", tC.id), "")

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_createStoreHandler(t *testing.T) {
	testCases := []struct {
		desc  string
		store model.Store
		err   error
		input string
		rcode int
		rdata string
	}{
		{
			desc:  "success",
			store: model.Store{Name: "Test1"},
			err:   nil,
			input: `{"name": "Test1"}`,
			rcode: http.StatusOK,
			rdata: `{"id":1, "name":"Test1"}`,
		},
		{
			desc:  "invalid: missing name",
			store: model.Store{},
			input: `{}`,
			err:   errSkip,
			rcode: http.StatusBadRequest,
			rdata: `{"error":"name is a required field"}`,
		},
		{
			desc:  "internal error",
			store: model.Store{Name: "Test1"},
			err:   errTest,
			input: `{"name": "Test1"}`,
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
				svc.EXPECT().CreateStore(gomock.Any(), tC.store).DoAndReturn(func(_ context.Context, s model.Store) (model.Store, error) {
					s.ID = 1
					return s, tC.err
				})
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodPost, "/v1/stores", tC.input)

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_updateStoreHandler(t *testing.T) {
	testCases := []struct {
		desc  string
		store model.Store
		err   error
		id    string
		input string
		rcode int
		rdata string
	}{
		{
			desc:  "success",
			store: model.Store{ID: 1, Name: "Test1"},
			err:   nil,
			id:    "1",
			input: `{"name": "Test1"}`,
			rcode: http.StatusOK,
			rdata: `{"id":1, "name":"Test1"}`,
		},
		{
			desc:  "invalid: missing name",
			store: model.Store{},
			id:    "1",
			input: `{}`,
			err:   errSkip,
			rcode: http.StatusBadRequest,
			rdata: `{"error":"name is a required field"}`,
		},
		{
			desc:  "invalid store ID",
			store: model.Store{},
			id:    "test",
			input: `{"name": "Test1"}`,
			err:   errSkip,
			rcode: http.StatusBadRequest,
			rdata: `{"error":"invalid store ID"}`,
		},
		{
			desc:  "not found",
			store: model.Store{ID: 1, Name: "Test1"},
			id:    "1",
			input: `{"name": "Test1"}`,
			err:   service.ErrNotFound,
			rcode: http.StatusNotFound,
			rdata: `{"error":"store not found"}`,
		},
		{
			desc:  "internal error",
			store: model.Store{ID: 1, Name: "Test1"},
			err:   errTest,
			id:    "1",
			input: `{"name": "Test1"}`,
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
				svc.EXPECT().UpdateStore(gomock.Any(), tC.store).Return(tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodPut, fmt.Sprintf("/v1/stores/%s", tC.id), tC.input)

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_deleteStoreHandler(t *testing.T) {
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
			desc:  "invalid store ID",
			id:    "test",
			err:   errSkip,
			rcode: http.StatusBadRequest,
			rdata: `{"error":"invalid store ID"}`,
		},
		{
			desc:  "not found",
			id:    "1",
			err:   service.ErrNotFound,
			rcode: http.StatusNotFound,
			rdata: `{"error":"store not found"}`,
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
				svc.EXPECT().DeleteStore(gomock.Any(), int64(1)).Return(tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodDelete, fmt.Sprintf("/v1/stores/%s", tC.id), "")

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

func Test_getStorePositionsHandler(t *testing.T) {
	testCases := []struct {
		desc      string
		storeID   string
		positions []model.Position
		err       error
		rcode     int
		rdata     string
	}{
		{
			desc:    "success",
			storeID: "1",
			positions: []model.Position{
				{ProductID: 1, StoreID: 1, Price: decimal.NewFromInt(100)},
				{ProductID: 2, StoreID: 1, Price: decimal.NewFromInt(200)},
			},
			err:   nil,
			rcode: http.StatusOK,
			rdata: `[{"productId":1, "storeId":1, "price":100},
				{"productId":2, "storeId":1, "price":200}]`,
		},
		{
			desc:      "internal error",
			storeID:   "1",
			positions: nil,
			err:       errTest,
			rcode:     http.StatusInternalServerError,
			rdata:     `{"error":"internal error"}`,
		},
		{
			desc:      "empty store ID",
			storeID:   "",
			positions: []model.Position{},
			err:       errSkip,
			rcode:     http.StatusBadRequest,
			rdata:     `{"error":"invalid store ID"}`,
		},
		{
			desc:      "invalid store ID",
			storeID:   "test",
			positions: []model.Position{},
			err:       errSkip,
			rcode:     http.StatusBadRequest,
			rdata:     `{"error":"invalid store ID"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().GetStorePositions(gomock.Any(), int64(1)).Return(tC.positions, tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodGet, fmt.Sprintf("/v1/stores/%s/positions", tC.storeID), "")

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_getCategoryProductsHandler(t *testing.T) {
	testCases := []struct {
		desc       string
		categoryID string
		products   []model.Product
		err        error
		rcode      int
		rdata      string
	}{
		{
			desc:       "success",
			categoryID: "1",
			products: []model.Product{
				{ID: 1, CategoryID: 1, Name: "Test1", Description: "Desc 1"},
				{ID: 2, CategoryID: 1, Name: "Test2", Description: "Desc 2"},
			},
			err:   nil,
			rcode: http.StatusOK,
			rdata: `[{"id":1, "categoryId":1, "name":"Test1", "description":"Desc 1"},
				{"id":2, "categoryId":1, "name":"Test2", "description":"Desc 2"}]`,
		},
		{
			desc:       "internal error",
			categoryID: "1",
			products:   nil,
			err:        errTest,
			rcode:      http.StatusInternalServerError,
			rdata:      `{"error":"internal error"}`,
		},
		{
			desc:       "empty store ID",
			categoryID: "",
			products:   []model.Product{},
			err:        errSkip,
			rcode:      http.StatusBadRequest,
			rdata:      `{"error":"invalid category ID"}`,
		},
		{
			desc:       "invalid store ID",
			categoryID: "test",
			products:   []model.Product{},
			err:        errSkip,
			rcode:      http.StatusBadRequest,
			rdata:      `{"error":"invalid category ID"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().GetProducts(gomock.Any(), int64(1)).Return(tC.products, tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodGet, fmt.Sprintf("/v1/categories/%s/products", tC.categoryID), "")

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}

func Test_getProductOffersHandler(t *testing.T) {
	testCases := []struct {
		desc      string
		productID string
		positions []model.Position
		err       error
		rcode     int
		rdata     string
	}{
		{
			desc:      "success",
			productID: "1",
			positions: []model.Position{
				{ProductID: 1, StoreID: 1, Price: decimal.NewFromInt(100)},
				{ProductID: 1, StoreID: 2, Price: decimal.NewFromInt(200)},
			},
			err:   nil,
			rcode: http.StatusOK,
			rdata: `[{"productId":1, "storeId":1, "price":100},
				{"productId":1, "storeId":2, "price":200}]`,
		},
		{
			desc:      "internal error",
			productID: "1",
			positions: nil,
			err:       errTest,
			rcode:     http.StatusInternalServerError,
			rdata:     `{"error":"internal error"}`,
		},
		{
			desc:      "empty product ID",
			productID: "",
			positions: []model.Position{},
			err:       errSkip,
			rcode:     http.StatusBadRequest,
			rdata:     `{"error":"invalid product ID"}`,
		},
		{
			desc:      "invalid product ID",
			productID: "test",
			positions: []model.Position{},
			err:       errSkip,
			rcode:     http.StatusBadRequest,
			rdata:     `{"error":"invalid product ID"}`,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			svc := service.NewMockService(ctrl)
			if tC.err != errSkip {
				svc.EXPECT().GetProductPositions(gomock.Any(), int64(1)).Return(tC.positions, tC.err)
			}

			router := setupTestRouter(svc)
			rec, r := newTestParameters(http.MethodGet, fmt.Sprintf("/v1/products/%s/offers", tC.productID), "")

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}
