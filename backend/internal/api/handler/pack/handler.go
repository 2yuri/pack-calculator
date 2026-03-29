package pack

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
	Service domain.PackService
}

type handler struct {
	service domain.PackService
}

func NewHandler(deps Deps) rest.Routes {
	h := &handler{
		service: deps.Service,
	}

	return rest.Routes{
		{
			Method:  http.MethodGet,
			Path:    "/products/:id/packs",
			Handler: h.getByProductID,
		},
		{
			Method:  http.MethodPost,
			Path:    "/products/:id/packs",
			Handler: h.create,
		},
		{
			Method:  http.MethodPost,
			Path:    "/products/:id/packs/batch",
			Handler: h.createBatch,
		},
		{
			Method:  http.MethodPut,
			Path:    "/packs/:id",
			Handler: h.update,
		},
		{
			Method:  http.MethodDelete,
			Path:    "/packs/:id",
			Handler: h.delete,
		},
	}
}

func (h *handler) getByProductID(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return rest.BadRequest(c, "invalid product id")
	}

	packs, err := h.service.GetByProductID(c.Request().Context(), productID)
	if err != nil {
		zap.L().Error("failed to get packs", zap.Error(err))
		return rest.InternalError(c)
	}

	if len(packs) == 0 {
		return rest.NoContent(c)
	}

	return rest.Ok(c, packs)
}

func (h *handler) create(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return rest.BadRequest(c, "invalid product id")
	}

	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return rest.BadRequest(c, "invalid request body")
	}

	if req.Size <= 0 {
		return rest.BadRequest(c, "size must be greater than 0")
	}

	pack, err := h.service.Create(c.Request().Context(), productID, req.Size)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicatePackSize) {
			return rest.Conflict(c, err.Error())
		}

		zap.L().Error("failed to create pack", zap.Error(err))
		return rest.InternalError(c)
	}

	return rest.Created(c, pack)
}

func (h *handler) createBatch(c echo.Context) error {
	productID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return rest.BadRequest(c, "invalid product id")
	}

	var req CreateBatchRequest
	if err := c.Bind(&req); err != nil {
		return rest.BadRequest(c, "invalid request body")
	}

	if len(req.Sizes) == 0 {
		return rest.BadRequest(c, "sizes must not be empty")
	}

	for _, size := range req.Sizes {
		if size <= 0 {
			return rest.BadRequest(c, "all sizes must be greater than 0")
		}
	}

	packs, err := h.service.CreateBatch(c.Request().Context(), productID, req.Sizes)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicatePackSize) {
			return rest.Conflict(c, err.Error())
		}

		zap.L().Error("failed to create packs batch", zap.Error(err))
		return rest.InternalError(c)
	}

	return rest.Created(c, packs)
}

func (h *handler) update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return rest.BadRequest(c, "invalid id")
	}

	var req CreateRequest
	if err := c.Bind(&req); err != nil {
		return rest.BadRequest(c, "invalid request body")
	}

	if req.Size <= 0 {
		return rest.BadRequest(c, "size must be greater than 0")
	}

	pack, err := h.service.Update(c.Request().Context(), id, req.Size)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return rest.NotFound(c, err.Error())
		}
		if errors.Is(err, domain.ErrDuplicatePackSize) {
			return rest.Conflict(c, err.Error())
		}

		zap.L().Error("failed to update pack", zap.Error(err))
		return rest.InternalError(c)
	}

	return rest.Ok(c, pack)
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

		zap.L().Error("failed to delete pack", zap.Error(err))
		return rest.InternalError(c)
	}

	return rest.NoContent(c)
}
