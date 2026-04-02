package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSubmitService_UpdatePackSizes(t *testing.T) {
	t.Run("updates pack sizes in storage", func(t *testing.T) {
		store := &mockPackStore{sizes: []int{250, 500}}
		svc := NewSubmitService(store)
		err := svc.UpdatePackSizes([]int{100, 200, 300})
		require.NoError(t, err)
		assert.Equal(t, []int{100, 200, 300}, store.sizes)
	})

	t.Run("storage error is propagated", func(t *testing.T) {
		svc := NewSubmitService(&mockPackStore{updateErr: errors.New("disk full")})
		err := svc.UpdatePackSizes([]int{100})
		assert.ErrorContains(t, err, "disk full")
	})
}