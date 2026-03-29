package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Route struct {
	Method  string
	Path    string
	Handler echo.HandlerFunc
}

type Routes []Route

func (r Routes) Register(g *echo.Group) {
	for _, route := range r {
		g.Add(route.Method, route.Path, route.Handler)
	}
}

type ErrResponse struct {
	Error string `json:"error,omitempty"`
}

func Ok(c echo.Context, data any) error {
	return c.JSON(http.StatusOK, data)
}

func Created(c echo.Context, data any) error {
	return c.JSON(http.StatusCreated, data)
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func BadRequest(c echo.Context, msg string) error {
	return c.JSON(http.StatusBadRequest, ErrResponse{Error: msg})
}

func NotFound(c echo.Context, msg string) error {
	return c.JSON(http.StatusNotFound, ErrResponse{Error: msg})
}

func Conflict(c echo.Context, msg string) error {
	return c.JSON(http.StatusConflict, ErrResponse{Error: msg})
}

func UnprocessableEntity(c echo.Context, msg string) error {
	return c.JSON(http.StatusUnprocessableEntity, ErrResponse{Error: msg})
}

func InternalError(c echo.Context) error {
	return c.JSON(http.StatusInternalServerError, ErrResponse{Error: "internal server error"})
}
