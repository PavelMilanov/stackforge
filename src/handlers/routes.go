package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

/*
registerPageRoutes описывает server-side rendered страницы приложения.

Входные параметры:
- e: Echo router, в который регистрируются маршруты страниц и HTMX-фрагментов.

Возвращает:
- ничего.
*/
func (h *Handler) registerPageRoutes(e *echo.Echo) {
	e.GET("/", h.templateCatalog)
	e.GET("/templates/preview", h.templatePreview)
	e.GET("/stands", h.stands)

	stands := e.Group("/stands")
	stands.GET("/create-modal", h.createStandModal)
	stands.GET("/stacks", h.standStacks)
	stands.POST("/create", h.createStand)
	stands.GET("/close-modal", h.closeStandModal)

	e.GET("/history", h.history)
	e.GET("/docs", h.docs)
}

/*
check возвращает минимальный healthcheck без обращения к внешним интеграциям.

Входные параметры:
- c: текущий Echo context HTTP-запроса.

Возвращает:
- error: ошибка JSON-ответа Echo или nil при успешной записи healthcheck.
*/
func (h *Handler) check(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
