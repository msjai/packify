package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHandler_Handle(t *testing.T) {
	tests := []struct {
		name     string
		mock     *mockGetter
		wantCode int
		wantBody string
	}{
		{
			name:     "returns pack sizes",
			mock:     &mockGetter{sizes: []int{250, 500, 1000}},
			wantCode: http.StatusOK,
			wantBody: `{"packs":[250,500,1000]}`,
		},
		{
			name:     "empty sizes",
			mock:     &mockGetter{sizes: []int{}},
			wantCode: http.StatusOK,
			wantBody: `{"packs":[]}`,
		},
		{
			name:     "nil sizes returns null",
			mock:     &mockGetter{},
			wantCode: http.StatusOK,
			wantBody: `{"packs":null}`,
		},
		{
			name:     "service error",
			mock:     &mockGetter{err: errors.New("db locked")},
			wantCode: http.StatusInternalServerError,
			wantBody: "failed to get pack sizes: db locked\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/packs", nil)
			w := httptest.NewRecorder()

			h := NewGetHandler(tc.mock)
			h.Handle(w, req)

			assert.Equal(t, tc.wantCode, w.Code)
			assert.Equal(t, tc.wantBody, w.Body.String())

			if tc.wantCode == http.StatusOK {
				assert.Equal(t, applicationJSON, w.Header().Get("Content-Type"))
			}
		})
	}
}
