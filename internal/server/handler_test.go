package server

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus/hooks/test"
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

			router := chi.NewRouter()
			SetupRouter(svc, router)

			test.NewGlobal()
			rec := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/v1/categories", nil)

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
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

			router := chi.NewRouter()
			SetupRouter(svc, router)

			test.NewGlobal()
			rec := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/v1/stores", nil)

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
			desc:    "empty storeId",
			storeID: "",
			items:   []*model.Item{},
			err:     errSkip,
			rcode:   http.StatusBadRequest,
			rdata:   `{"error":"invalid storeId"}`,
		},
		{
			desc:    "invalid storeId",
			storeID: "test",
			items:   []*model.Item{},
			err:     errSkip,
			rcode:   http.StatusBadRequest,
			rdata:   `{"error":"invalid storeId"}`,
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

			router := chi.NewRouter()
			SetupRouter(svc, router)

			test.NewGlobal()
			rec := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/stores/%s/items", tC.storeID), nil)

			router.ServeHTTP(rec, r)

			body, _ := ioutil.ReadAll(rec.Result().Body)

			assert.Equal(t, tC.rcode, rec.Result().StatusCode)
			assert.JSONEq(t, tC.rdata, string(body))
		})
	}
}
