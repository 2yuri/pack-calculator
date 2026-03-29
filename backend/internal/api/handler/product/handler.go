package product

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"

	"github.com/2yuri/pack-calculator/internal/domain"
	"github.com/2yuri/pack-calculator/pkg/rest"
)

type Deps struct {
	Service domain.ProductService
}

type handler struct {
	service domain.ProductService
}

func NewHandler(deps Deps) rest.Routes {
	h := &handler{
		service: deps.Service,
	}

	return rest.Routes{
		{
			Method:  http.MethodGet,
			Path:    "/products",
			Handler: h.getAll,
		},
		{
			Method:  http.MethodGet,
			Path:    "/products/:id",
			Handler: h.getByID,
		},
		{Method: http.MethodPost,
			Path:    "/products",
			Handler: h.create,
		},

		{Method: http.MethodPut,
			Path:    "/products/:id",
			Handler: h.update,
		},

		{Method: http.MethodDelete,
			Path:    "/products/:id",
			Handler: h.delete,
		},
	}
}

func (h *handler) getAll(c echo.Context) error {
	products, err := h.service.GetAll(c.Request().Context())
	if err != nil {
		zap.L().Error("failed to get products", zap.Error(err))
		return rest.InternalError(c)
	}

	if len(products) == 0 {
		return rest.NoContent(c)
	}

	return rest.Ok(c, products)
}

func (h *handler) getByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return rest.BadRequest(c, "invalid id")
	}

	product, err := h.service.GetByID(c.Request().Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return rest.NotFound(c, err.Error())
		}

		zap.L().Error("failed to get product", zap.Error(err))
		return rest.InternalError(c)
	}

	return rest.Ok(c, product)
}

func (h *handler) create(c echo.Context) error {
	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return rest.BadRequest(c, "invalid request body")
	}

	if req.Name == "" {
		return rest.BadRequest(c, "name is required")
	}

	product, err := h.service.Create(c.Request().Context(), req.Name)
	if err != nil {
		zap.L().Error("failed to create product", zap.Error(err))
		return rest.InternalError(c)
	}

	return rest.Created(c, product)
}

func (h *handler) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return rest.BadRequest(c, "invalid id")
	}

	var req UpdateRequest
	if err := c.Bind(&req); err != nil {
		return rest.BadRequest(c, "invalid request body")
	}

	if req.Name == "" {
		return rest.BadRequest(c, "name is required")
	}

	product, err := h.service.Update(c.Request().Context(), id, req.Name)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return rest.NotFound(c, err.Error())
		}

		zap.L().Error("failed to update product", zap.Error(err))
		return rest.InternalError(c)
	}

	return rest.Ok(c, product)
}

func (h *handler) delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return rest.BadRequest(c, "invalid id")
	}

	if err := h.service.Delete(c.Request().Context(), id); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return rest.NotFound(c, err.Error())
		}

		zap.L().Error("failed to delete product", zap.Error(err))
		return rest.InternalError(c)
	}
	return rest.NoContent(c)
}
