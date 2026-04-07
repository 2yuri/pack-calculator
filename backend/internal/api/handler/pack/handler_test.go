package pack

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

func TestHandler_getByProductID(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		mockSetup    func(svc *mocks.MockPackService)
		expectedCode int
		expectedBody string
	}{
		{
			name:    "success",
			paramID: "1",
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().GetByProductID(gomock.Any(), 1).Return([]*domain.Pack{
					{ID: 1, ProductID: 1, Size: 250},
				}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"size":250`,
		},
		{
			name:    "empty list",
			paramID: "1",
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().GetByProductID(gomock.Any(), 1).Return([]*domain.Pack{}, nil)
			},
			expectedCode: http.StatusNoContent,
		},
		{
			name:         "invalid product id",
			paramID:      "abc",
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid product id",
		},
		{
			name:    "internal error",
			paramID: "1",
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().GetByProductID(gomock.Any(), 1).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockPackService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/products/"+tt.paramID+"/packs", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			err := h.getByProductID(c)
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
		paramID      string
		body         string
		mockSetup    func(svc *mocks.MockPackService)
		expectedCode int
		expectedBody string
	}{
		{
			name:    "success",
			paramID: "1",
			body:    `{"size":250}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().Create(gomock.Any(), 1, 250).Return(&domain.Pack{ID: 1, ProductID: 1, Size: 250}, nil)
			},
			expectedCode: http.StatusCreated,
			expectedBody: `"size":250`,
		},
		{
			name:         "invalid product id",
			paramID:      "abc",
			body:         `{"size":250}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid product id",
		},
		{
			name:         "invalid body",
			paramID:      "1",
			body:         `{invalid`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid request body",
		},
		{
			name:         "size zero",
			paramID:      "1",
			body:         `{"size":0}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "size must be greater than 0",
		},
		{
			name:         "negative size",
			paramID:      "1",
			body:         `{"size":-1}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "size must be greater than 0",
		},
		{
			name:    "duplicate pack size",
			paramID: "1",
			body:    `{"size":250}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().Create(gomock.Any(), 1, 250).Return(nil, domain.ErrDuplicatePackSize)
			},
			expectedCode: http.StatusConflict,
			expectedBody: domain.ErrDuplicatePackSize.Error(),
		},
		{
			name:    "internal error",
			paramID: "1",
			body:    `{"size":250}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().Create(gomock.Any(), 1, 250).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockPackService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/products/"+tt.paramID+"/packs", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			err := h.create(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}

func TestHandler_createBatch(t *testing.T) {
	tests := []struct {
		name         string
		paramID      string
		body         string
		mockSetup    func(svc *mocks.MockPackService)
		expectedCode int
		expectedBody string
	}{
		{
			name:    "success",
			paramID: "1",
			body:    `{"sizes":[250,500]}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().CreateBatch(gomock.Any(), 1, []int{250, 500}).Return([]*domain.Pack{
					{ID: 1, ProductID: 1, Size: 250},
					{ID: 2, ProductID: 1, Size: 500},
				}, nil)
			},
			expectedCode: http.StatusCreated,
		},
		{
			name:         "invalid product id",
			paramID:      "abc",
			body:         `{"sizes":[250]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid product id",
		},
		{
			name:         "invalid body",
			paramID:      "1",
			body:         `{invalid`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid request body",
		},
		{
			name:         "empty sizes",
			paramID:      "1",
			body:         `{"sizes":[]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "sizes must not be empty",
		},
		{
			name:         "size zero in batch",
			paramID:      "1",
			body:         `{"sizes":[250,0]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "all sizes must be greater than 0",
		},
		{
			name:    "duplicate pack size",
			paramID: "1",
			body:    `{"sizes":[250]}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().CreateBatch(gomock.Any(), 1, []int{250}).Return(nil, domain.ErrDuplicatePackSize)
			},
			expectedCode: http.StatusConflict,
		},
		{
			name:    "internal error",
			paramID: "1",
			body:    `{"sizes":[250]}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().CreateBatch(gomock.Any(), 1, []int{250}).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockPackService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/products/"+tt.paramID+"/packs/batch", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(tt.paramID)

			err := h.createBatch(c)
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
		mockSetup    func(svc *mocks.MockPackService)
		expectedCode int
		expectedBody string
	}{
		{
			name:    "success",
			paramID: "1",
			body:    `{"size":500}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().Update(gomock.Any(), 1, 500).Return(&domain.Pack{ID: 1, ProductID: 1, Size: 500}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"size":500`,
		},
		{
			name:         "invalid id",
			paramID:      "abc",
			body:         `{"size":500}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid id",
		},
		{
			name:         "invalid body",
			paramID:      "1",
			body:         `{invalid`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid request body",
		},
		{
			name:         "size zero",
			paramID:      "1",
			body:         `{"size":0}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "size must be greater than 0",
		},
		{
			name:    "not found",
			paramID: "999",
			body:    `{"size":500}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().Update(gomock.Any(), 999, 500).Return(nil, domain.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "duplicate pack size",
			paramID: "1",
			body:    `{"size":500}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().Update(gomock.Any(), 1, 500).Return(nil, domain.ErrDuplicatePackSize)
			},
			expectedCode: http.StatusConflict,
		},
		{
			name:    "internal error",
			paramID: "1",
			body:    `{"size":500}`,
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().Update(gomock.Any(), 1, 500).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockPackService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPut, "/packs/"+tt.paramID, strings.NewReader(tt.body))
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
		mockSetup    func(svc *mocks.MockPackService)
		expectedCode int
	}{
		{
			name:    "success",
			paramID: "1",
			mockSetup: func(svc *mocks.MockPackService) {
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
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().Delete(gomock.Any(), 999).Return(domain.ErrNotFound)
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name:    "internal error",
			paramID: "1",
			mockSetup: func(svc *mocks.MockPackService) {
				svc.EXPECT().Delete(gomock.Any(), 1).Return(errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockPackService(ctrl)
			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}
			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodDelete, "/packs/"+tt.paramID, nil)
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
