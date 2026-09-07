package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cyntraten/metrics-ed/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestHandlerUpdateMetrics(t *testing.T) {
	tests := []struct {
		name               string
		method             string
		url                string
		expectedStatusCode int
	}{
		{
			name:               "Test update gauge metric",
			method:             http.MethodPost,
			url:                "/update/gauge/test/10.5",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Test update counter metric",
			method:             http.MethodPost,
			url:                "/update/counter/test/12",
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Test update counter metric error Bad Request",
			method:             http.MethodPost,
			url:                "/update/counter/test/12.5",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Test update gauge metric error Bad Request",
			method:             http.MethodPost,
			url:                "/update/gauge/test/ggrgah",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name:               "Test update gauge metric error Not Found",
			method:             http.MethodPost,
			url:                "/update/gauge//ggrgah",
			expectedStatusCode: http.StatusNotFound,
		},
		{
			name:               "Test not allowed method",
			method:             http.MethodGet,
			url:                "/update/gauge//ggrgah",
			expectedStatusCode: http.StatusMethodNotAllowed,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := repository.NewMemStorage()
			h := NewHandler(storage)
			req := httptest.NewRequest(test.method, test.url, nil)
			w := httptest.NewRecorder()

			h.UpdateMetrics(w, req)

			res := w.Result()
			defer res.Body.Close()

			assert.Equal(t, res.StatusCode, test.expectedStatusCode)
		})

	}
}
