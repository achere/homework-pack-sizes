package server

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type SizeRepoStub struct {
	getPackSizes  func(ctx context.Context) ([]int, error)
	storePackSies func(ctx context.Context, sizes []int) error
}

func (sr *SizeRepoStub) GetPackSizes(ctx context.Context) ([]int, error) {
	return sr.getPackSizes(ctx)
}

func (sr *SizeRepoStub) StorePackSizes(ctx context.Context, sizes []int) error {
	return sr.storePackSies(ctx, sizes)
}

func TestCalculatePacksHandler(t *testing.T) {
	app := NewTestApp()
	app.SizeRepo = &SizeRepoStub{
		getPackSizes: func(ctx context.Context) ([]int, error) {
			return []int{250, 500, 1000, 2000, 5000}, nil
		},
	}

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedPacks  map[int]int
		expectedError  bool
	}{
		{
			name:           "Valid request",
			requestBody:    `{"order": 12500}`,
			expectedStatus: http.StatusOK,
			expectedPacks:  map[int]int{5000: 2, 2000: 1, 500: 1},
			expectedError:  false,
		},
		{
			name:           "Negative order",
			requestBody:    `{"order": -12500}`,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"order": "251"}`,
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
			req := httptest.NewRequest(http.MethodPost, "/api/v2/calculate-packs", body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			app.calculatePacksHandler(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)

			var resp calculatePacksResponseV1
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

func TestCalculatePacksHandlerV1(t *testing.T) {
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

			app.calculatePacksHandlerV1(rr, req)

			assert.Equal(t, test.expectedStatus, rr.Code)

			var resp calculatePacksResponseV1
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
	app.SizeRepo = &SizeRepoStub{
		getPackSizes: func(ctx context.Context) ([]int, error) {
			return []int{250, 500, 1000, 2000, 5000}, nil
		},
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.uiHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "<title>Pack Calculator</title>")
}

func TestRetrievePackSizesHandler(t *testing.T) {
	tests := []struct {
		name           string
		getPackSizes   func(ctx context.Context) ([]int, error)
		expectedStatus int
		expectedSizes  []int
		expectedError  bool
	}{
		{
			name: "Success",
			getPackSizes: func(ctx context.Context) ([]int, error) {
				return []int{250, 500, 1000}, nil
			},
			expectedStatus: http.StatusOK,
			expectedSizes:  []int{250, 500, 1000},
			expectedError:  false,
		},
		{
			name: "Error from repo",
			getPackSizes: func(ctx context.Context) ([]int, error) {
				return nil, assert.AnError
			},
			expectedStatus: http.StatusInternalServerError,
			expectedSizes:  nil,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewTestApp()
			app.SizeRepo = &SizeRepoStub{
				getPackSizes: tt.getPackSizes,
			}

			req := httptest.NewRequest(http.MethodGet, "/api/v2/sizes", nil)
			rr := httptest.NewRecorder()

			app.retrievePackSizesHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			var resp retrievePackSizesResponse
			err := json.Unmarshal(rr.Body.Bytes(), &resp)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.NotEmpty(t, resp.Error)
				assert.Nil(t, resp.Sizes)
			} else {
				assert.Empty(t, resp.Error)
				assert.Equal(t, tt.expectedSizes, resp.Sizes)
			}
		})
	}
}

func TestStorePackSizesHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    string
		storePackSies  func(ctx context.Context, sizes []int) error
		expectedStatus int
		expectError    bool
	}{
		{
			name:        "Success",
			requestBody: `{"sizes": [250, 500, 1000]}`,
			storePackSies: func(ctx context.Context, sizes []int) error {
				assert.Equal(t, []int{250, 500, 1000}, sizes)
				return nil
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:        "Invalid JSON",
			requestBody: `{"sizes": ["250", 500, 1000]}`,
			storePackSies: func(ctx context.Context, sizes []int) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:        "Empty body",
			requestBody: ``,
			storePackSies: func(ctx context.Context, sizes []int) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:        "Error from repo",
			requestBody: `{"sizes": [250, 500, 1000]}`,
			storePackSies: func(ctx context.Context, sizes []int) error {
				return assert.AnError
			},
			expectedStatus: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := NewTestApp()
			app.SizeRepo = &SizeRepoStub{
				storePackSies: tt.storePackSies,
			}

			body := bytes.NewBufferString(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v2/sizes", body)
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()

			app.storePackSizesHandler(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectError {
				var resp storePackSizesResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Error)
			}
		})
	}
}

func NewTestApp() *App {
	return &App{
		logger: slog.New(slog.DiscardHandler),
		Config: &Config{
			Order: "250",
		},
	}
}
