package handlers

import (
	"net/http"

	"github.com/PavelMilanov/stackforge/views/pages"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
)

func (h *Handler) registerPageRoutes(e *echo.Echo) {
	e.GET("/", h.templateCatalog)
	e.GET("/stands", h.stands)
	e.GET("/history", h.history)
	e.GET("/docs", h.docs)
}

func (h *Handler) render(c *echo.Context, component templ.Component) error {
	c.Response().WriteHeader(http.StatusOK)
	return component.Render(c.Request().Context(), c.Response())
}

func (h *Handler) templateCatalog(c *echo.Context) error {
	return h.render(c, pages.TemplateCatalogPage())
}

func (h *Handler) stands(c *echo.Context) error {
	return h.render(c, pages.StandsPage())
}

func (h *Handler) history(c *echo.Context) error {
	return h.render(c, pages.HistoryPage())
}

func (h *Handler) docs(c *echo.Context) error {
	return h.render(c, pages.DocsPage())
}
