package order

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

func TestHandler_calculate(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		mockSetup    func(svc *mocks.MockOrderService)
		expectedCode int
		expectedBody string
	}{
		{
			name: "success",
			body: `{"items":[{"product_id":1,"quantity":1}]}`,
			mockSetup: func(svc *mocks.MockOrderService) {
				svc.EXPECT().
					Calculate(gomock.Any(), []domain.OrderItem{{ProductID: 1, Quantity: 1}}).
					Return([]*domain.CalculationResult{
						{Quantity: 1, ProductID: 1, TotalItems: 250, TotalPacks: 1},
					}, nil)
			},
			expectedCode: http.StatusOK,
			expectedBody: `"results"`,
		},
		{
			name:         "empty items",
			body:         `{"items":[]}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "items must not be empty",
		},
		{
			name:         "invalid body",
			body:         `{invalid`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "invalid request body",
		},
		{
			name: "invalid quantity",
			body: `{"items":[{"product_id":1,"quantity":0}]}`,
			mockSetup: func(svc *mocks.MockOrderService) {
				svc.EXPECT().Calculate(gomock.Any(), gomock.Any()).Return(nil, domain.ErrInvalidQuantity)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: domain.ErrInvalidQuantity.Error(),
		},
		{
			name: "duplicate product id",
			body: `{"items":[{"product_id":1,"quantity":1},{"product_id":1,"quantity":2}]}`,
			mockSetup: func(svc *mocks.MockOrderService) {
				svc.EXPECT().Calculate(gomock.Any(), gomock.Any()).Return(nil, domain.ErrDuplicateProductID)
			},
			expectedCode: http.StatusBadRequest,
			expectedBody: domain.ErrDuplicateProductID.Error(),
		},
		{
			name: "no packs available",
			body: `{"items":[{"product_id":1,"quantity":1}]}`,
			mockSetup: func(svc *mocks.MockOrderService) {
				svc.EXPECT().Calculate(gomock.Any(), gomock.Any()).Return(nil, domain.ErrNoPacksAvailable)
			},
			expectedCode: http.StatusUnprocessableEntity,
			expectedBody: domain.ErrNoPacksAvailable.Error(),
		},
		{
			name: "internal error",
			body: `{"items":[{"product_id":1,"quantity":1}]}`,
			mockSetup: func(svc *mocks.MockOrderService) {
				svc.EXPECT().Calculate(gomock.Any(), gomock.Any()).Return(nil, errors.New("db error"))
			},
			expectedCode: http.StatusInternalServerError,
			expectedBody: "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			svc := mocks.NewMockOrderService(ctrl)

			if tt.mockSetup != nil {
				tt.mockSetup(svc)
			}

			h := &handler{service: svc}

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err := h.calculate(c)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, rec.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, rec.Body.String(), tt.expectedBody)
			}
		})
	}
}
