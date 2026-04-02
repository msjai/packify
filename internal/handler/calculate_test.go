package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     CalculateRequest
		wantErr string
	}{
		{
			name:    "valid order",
			req:     CalculateRequest{Order: 500},
			wantErr: "",
		},
		{
			name:    "zero order",
			req:     CalculateRequest{Order: 0},
			wantErr: "order must be positive",
		},
		{
			name:    "negative order",
			req:     CalculateRequest{Order: -10},
			wantErr: "order must be positive",
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