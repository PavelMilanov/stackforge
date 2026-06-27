package handlers

import (
	"net/http"

	"github.com/labstack/echo/v5"
)

/*
check возвращает простой healthcheck-ответ.

Параметры:
  - c: Echo context входящего HTTP-запроса.

Результат: HTTP 200 с JSON {"status":"ok"} или ошибка записи ответа.
*/
func (h *Handler) check(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
