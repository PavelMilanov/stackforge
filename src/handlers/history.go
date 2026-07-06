package handlers

import (
	"github.com/PavelMilanov/stackforge/views/pages"
	"github.com/labstack/echo/v5"
)

/*
history отображает журнал запусков и операций со стендами.

Входные параметры:
- c: текущий Echo context HTTP-запроса.

Возвращает:
- error: ошибка рендера страницы или nil при успешной записи ответа.
*/
func (h *Handler) history(c *echo.Context) error {
	return h.render(c, pages.HistoryPage())
}
