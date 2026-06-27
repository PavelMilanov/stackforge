package handlers

import (
	"net/http"

	"github.com/PavelMilanov/stackforge/views/pages"
	"github.com/labstack/echo/v5"
)

func (h *Handler) registerPageRoutes(e *echo.Echo) {
	e.GET("/", h.dashboard)
	// e.GET("/templates", h.templateCatalog)
}

func (h *Handler) dashboard(c *echo.Context) error {
	component := pages.Dashboard()
	c.Response().WriteHeader(http.StatusOK)
	return component.Render(c.Request().Context(), c.Response())
}

// func (h *Handler) templateCatalog(c echo.Context) error {
// 	items := []pages.TemplateListItem{
// 		{
// 			ID:          "backend-basic",
// 			Name:        "Backend Basic",
// 			Description: "Базовый backend stack для разработки.",
// 		},
// 	}

// 	component := pages.TemplateCatalog(items)
// 	c.Response().WriteHeader(http.StatusOK)
// 	return component.Render(c.Request().Context(), c.Response().Writer)
// }
