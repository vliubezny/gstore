package server

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func Test_getLogger(t *testing.T) {
	l := logrus.New()
	ctx := context.WithValue(context.Background(), loggerKey{}, l)
	r := httptest.NewRequest("", "/", nil).WithContext(ctx)

	logger := getLogger(r)

	assert.Exactly(t, l, logger)
}
