package server

import (
	"errors"
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
			rdata: `{"categories":[{"id":1, "name":"Test1"}, {"id":2, "name":"Test2"}]}`,
		},
		{
			desc:       "internal error",
			categories: nil,
			err:        errors.New("test error"),
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
