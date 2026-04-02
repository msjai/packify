package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubmitHandler_Handle(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		mock      *mockSubmitter
		wantCode  int
		wantBody  string
		wantSizes []int
	}{
		{
			name:      "valid request",
			body:      `{"packs":[100,200,300]}`,
			mock:      &mockSubmitter{},
			wantCode:  http.StatusOK,
			wantSizes: []int{100, 200, 300},
		},
		{
			name:     "empty body",
			body:     ``,
			mock:     &mockSubmitter{},
			wantCode: http.StatusBadRequest,
			wantBody: "decode request: EOF\n",
		},
		{
			name:     "invalid JSON",
			body:     `{broken`,
			mock:     &mockSubmitter{},
			wantCode: http.StatusBadRequest,
			wantBody: "decode request: invalid character 'b' looking for beginning of object key string\n",
		},
		{
			name:     "empty pack list",
			body:     `{"packs":[]}`,
			mock:     &mockSubmitter{},
			wantCode: http.StatusBadRequest,
			wantBody: "pack sizes list is empty\n",
		},
		{
			name:     "duplicate sizes",
			body:     `{"packs":[250,250]}`,
			mock:     &mockSubmitter{},
			wantCode: http.StatusBadRequest,
			wantBody: "duplicate pack size\n",
		},
		{
			name:     "service error",
			body:     `{"packs":[100]}`,
			mock:     &mockSubmitter{err: errors.New("db write failed")},
			wantCode: http.StatusInternalServerError,
			wantBody: "failed to update pack sizes: db write failed\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/packs", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			h := NewSubmitHandler(tc.mock)
			h.Handle(w, req)

			assert.Equal(t, tc.wantCode, w.Code)

			if tc.wantBody != "" {
				assert.Equal(t, tc.wantBody, w.Body.String())
			}

			if tc.wantSizes != nil {
				assert.Equal(t, tc.wantSizes, tc.mock.sizes)
			}
		})
	}
}

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
