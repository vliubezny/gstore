package mem

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCategories(t *testing.T) {
	s := New()

	categories, err := s.GetCategories(context.Background())

	assert.NoError(t, err)
	assert.NotEmpty(t, categories, "Category slice is not populated")
}
