package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetService_GetPackSizes(t *testing.T) {
	t.Run("returns pack sizes from storage", func(t *testing.T) {
		svc := NewGetService(&mockPackStore{sizes: []int{250, 500, 1000}})
		sizes, err := svc.GetPackSizes()
		require.NoError(t, err)
		assert.Equal(t, []int{250, 500, 1000}, sizes)
	})

	t.Run("storage error is propagated", func(t *testing.T) {
		svc := NewGetService(&mockPackStore{getErr: errors.New("db locked")})
		_, err := svc.GetPackSizes()
		assert.ErrorContains(t, err, "db locked")
	})
}
