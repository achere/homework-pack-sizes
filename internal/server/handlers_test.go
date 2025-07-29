package server

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculatePacksHandler_Success(t *testing.T) {
	app := NewTestApp()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedPacks  map[int]int
		expectedError  bool
	}{
		{
			name:           "Valid request",
			requestBody:    `{"sizes": [250, 500, 1000, 2000, 5000], "order": 12500}`,
			expectedStatus: http.StatusOK,
			expectedPacks:  map[int]int{5000: 2, 2000: 1, 500: 1},
			expectedError:  false,
		},
		{
			name:           "Negative order",
			requestBody:    `{"sizes": [250, 500, 1000, 2000, 5000], "order": -12500}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "Negative size",
			requestBody:    `{"sizes": [-250, 500, 1000, 2000, 5000], "order": -12500}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"sizes": [250], "order": "251"}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "Empty body",
			requestBody:    ``,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			body := bytes.NewBufferString(test.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/calculate-packs", body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			app.calculatePacksHandler(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)

			var resp calculatePacksResponse
			err := json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if test.expectedError {
				assert.NotEmpty(t, resp.Error)
				assert.Empty(t, resp.Packs)
			} else {

				assert.Empty(t, resp.Error)
				assert.Equal(t, test.expectedPacks, resp.Packs)
			}
		})
	}
}

func TestUIHandler_Success(t *testing.T) {
	app := NewTestApp()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.uiHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "<title>Pack Calculator</title>")
}

func NewTestApp() *App {
	return &App{
		logger: slog.New(slog.DiscardHandler),
		Config: &Config{
			Sizes: []int{250, 500, 1000, 2000, 5000},
			Order: "250",
		},
	}
}
