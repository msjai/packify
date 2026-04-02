package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdatePacksRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     UpdatePacksRequest
		wantErr string
	}{
		{
			name:    "valid request",
			req:     UpdatePacksRequest{Packs: []int{250, 500, 1000}},
			wantErr: "",
		},
		{
			name:    "single pack size",
			req:     UpdatePacksRequest{Packs: []int{100}},
			wantErr: "",
		},
		{
			name:    "empty list",
			req:     UpdatePacksRequest{Packs: []int{}},
			wantErr: "pack sizes list is empty",
		},
		{
			name:    "nil list",
			req:     UpdatePacksRequest{Packs: nil},
			wantErr: "pack sizes list is empty",
		},
		{
			name:    "zero value",
			req:     UpdatePacksRequest{Packs: []int{250, 0, 500}},
			wantErr: "pack size must be positive",
		},
		{
			name:    "negative value",
			req:     UpdatePacksRequest{Packs: []int{250, -10, 500}},
			wantErr: "pack size must be positive",
		},
		{
			name:    "duplicate values",
			req:     UpdatePacksRequest{Packs: []int{250, 500, 250}},
			wantErr: "duplicate pack size",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.req.Validate()
			if tc.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tc.wantErr)
			}
		})
	}
}
