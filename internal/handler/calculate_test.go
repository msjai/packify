package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateHandler_Handle(t *testing.T) {
	tests := []struct {
		name      string
		body      string
		mock      *mockCalculator
		wantCode  int
		wantBody  string
		wantOrder int
	}{
		{
			name:      "valid request",
			body:      `{"order":501}`,
			mock:      &mockCalculator{result: map[int]int{500: 1, 250: 1}},
			wantCode:  http.StatusOK,
			wantBody:  `{"packs":[{"pack":250,"quantity":1},{"pack":500,"quantity":1}]}`,
			wantOrder: 501,
		},
		{
			name:     "empty body",
			body:     ``,
			mock:     &mockCalculator{},
			wantCode: http.StatusBadRequest,
			wantBody: "decode request: EOF\n",
		},
		{
			name:     "invalid JSON",
			body:     `{broken`,
			mock:     &mockCalculator{},
			wantCode: http.StatusBadRequest,
			wantBody: "decode request: invalid character 'b' looking for beginning of object key string\n",
		},
		{
			name:     "zero order",
			body:     `{"order":0}`,
			mock:     &mockCalculator{},
			wantCode: http.StatusBadRequest,
			wantBody: "order must be positive\n",
		},
		{
			name:     "negative order",
			body:     `{"order":-5}`,
			mock:     &mockCalculator{},
			wantCode: http.StatusBadRequest,
			wantBody: "order must be positive\n",
		},
		{
			name:     "service error",
			body:     `{"order":100}`,
			mock:     &mockCalculator{err: errors.New("no pack sizes configured")},
			wantCode: http.StatusInternalServerError,
			wantBody: "failed to calculate packs: no pack sizes configured\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/calculate", strings.NewReader(tc.body))
			w := httptest.NewRecorder()

			h := NewCalculateHandler(tc.mock)
			h.Handle(w, req)

			assert.Equal(t, tc.wantCode, w.Code)
			assert.Equal(t, tc.wantBody, w.Body.String())

			if tc.wantOrder != 0 {
				assert.Equal(t, tc.wantOrder, tc.mock.gotOrder)
			}
		})
	}
}

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