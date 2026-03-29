package order

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/2yuri/pack-calculator/pkg/rest"
)

type Deps struct {
	Service domain.OrderService
}

type handler struct {
	service domain.OrderService
}

func NewHandler(deps Deps) rest.Routes {
	h := &handler{
		service: deps.Service,
	}

	return rest.Routes{
		{
			Method:  http.MethodPost,
			Path:    "/orders",
			Handler: h.calculate,
		},
	}
}

func (h *handler) calculate(c echo.Context) error {
	var req CalculateRequest
	if err := c.Bind(&req); err != nil {
		return rest.BadRequest(c, "invalid request body")
	}

	if len(req.Items) == 0 {
		return rest.BadRequest(c, "items must not be empty")
	}

	items := make([]domain.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = domain.OrderItem{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
	}

	results, err := h.service.Calculate(c.Request().Context(), items)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidQuantity) || errors.Is(err, domain.ErrDuplicateProductID) {
			return rest.BadRequest(c, err.Error())
		}
		if errors.Is(err, domain.ErrNoPacksAvailable) {
			return rest.UnprocessableEntity(c, err.Error())
		}

		zap.L().Error("failed to calculate order", zap.Error(err))
		return rest.InternalError(c)
	}

	return rest.Ok(c, CalculateResponse{Results: results})
}
