package handlers

import (
	"github.com/PavelMilanov/stackforge/views/pages"
	"github.com/labstack/echo/v5"
)

/*
docs отображает короткие карточки документации со ссылками на внутреннюю wiki.

Входные параметры:
- c: HTTP-запрос.

Возвращает:
- error: ошибка рендера страницы или nil при успешной записи ответа.
*/
func (h *Handler) docs(c *echo.Context) error {
	return h.render(c, pages.DocsPage())
}
