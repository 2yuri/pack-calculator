package product

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/2yuri/pack-calculator/internal/mocks"
)

func TestHandler_getAll(t *testing.T) {
	tests := []struct {
		name         string
		mockSetup    func(svc *mocks.MockProductService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetAll(gomock.Any()).Return([]*domain.Product{
					{ID: 1, Name: "Product A"},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "Product A",
		},
		{
			name: "empty list",
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetAll(gomock.Any()).Return([]*domain.Product{}, nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name: "internal error",
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetAll(gomock.Any()).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockProductService(ctrl)
			tt.mockSetup(svc)
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/products", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.getAll(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHandler_getByID(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		mockSetup    func(svc *mocks.MockProductService)
		expectedCode int
		expectedBody string
	}{
		{
			name:    "success",
			paramID: "1",
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetByID(gomock.Any(), 1).Return(&domain.Product{ID: 1, Name: "Product A"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "Product A",
		},
		{
			name:         "invalid id",
			paramID:      "abc",
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid id",
		},
		{
			name:    "not found",
			paramID: "999",
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetByID(gomock.Any(), 999).Return(nil, domain.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "internal error",
			paramID: "1",
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().GetByID(gomock.Any(), 1).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockProductService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.paramID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			err := h.getByID(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHandler_create(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		mockSetup    func(svc *mocks.MockProductService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			body: `{"name":"New Product"}`,
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().Create(gomock.Any(), "New Product").Return(&domain.Product{ID: 1, Name: "New Product"}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: "New Product",
		},
		{
			name:         "empty name",
			body:         `{"name":""}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "name is required",
		},
		{
			name:         "invalid body",
			body:         `{invalid`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid request body",
		},
		{
			name: "internal error",
			body: `{"name":"Test"}`,
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().Create(gomock.Any(), "Test").Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockProductService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.create(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHandler_update(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		body         string
		mockSetup    func(svc *mocks.MockProductService)
		expectedCode int
		expectedBody string
	}{
		{
			name:    "success",
			paramID: "1",
			body:    `{"name":"Updated"}`,
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().Update(gomock.Any(), 1, "Updated").Return(&domain.Product{ID: 1, Name: "Updated"}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: "Updated",
		},
		{
			name:         "invalid id",
			paramID:      "abc",
			body:         `{"name":"X"}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid id",
		},
		{
			name:         "empty name",
			paramID:      "1",
			body:         `{"name":""}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "name is required",
		},
		{
			name:    "not found",
			paramID: "1",
			body:    `{"name":"X"}`,
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().Update(gomock.Any(), 1, "X").Return(nil, domain.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "internal error",
			paramID: "1",
			body:    `{"name":"X"}`,
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().Update(gomock.Any(), 1, "X").Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockProductService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/products/"+tt.paramID, strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			err := h.update(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHandler_delete(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		mockSetup    func(svc *mocks.MockProductService)
		expectedCode int
	}{
		{
			name:    "success",
			paramID: "1",
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().Delete(gomock.Any(), 1).Return(nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "invalid id",
			paramID:      "abc",
			expectedCode: http.StatusBadRequest,
		},
		{
			name:    "not found",
			paramID: "999",
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().Delete(gomock.Any(), 999).Return(domain.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "internal error",
			paramID: "1",
			mockSetup: func(svc *mocks.MockProductService) {
				svc.EXPECT().Delete(gomock.Any(), 1).Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockProductService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/products/"+tt.paramID, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			err := h.delete(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
		})
	}
}
