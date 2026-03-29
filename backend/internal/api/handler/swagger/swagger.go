package swagger

import (
	_ "embed"
	"net/http"

	"github.com/2yuri/pack-calculator/pkg/rest"
	"github.com/labstack/echo/v4"
)

//go:embed openapi.spec.yml
var openAPISpec []byte

//go:embed swagger.html
var swaggerHTML []byte

type handler struct {
}

func (h *handler) openapi(c echo.Context) error {
	return c.Blob(http.StatusOK, "application/yaml", openAPISpec)
}

func (h *handler) swagger(c echo.Context) error {
	return c.HTMLBlob(http.StatusOK, swaggerHTML)
}

func NewHandler() rest.Routes {
	h := &handler{}

	return rest.Routes{
		{Method: http.MethodGet,
			Path:    "/docs/openapi.yml",
			Handler: h.openapi,
		},

		{Method: http.MethodGet,
			Path:    "/docs",
			Handler: h.swagger,
		},
	}
}
