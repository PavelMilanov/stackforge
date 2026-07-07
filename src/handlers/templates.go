package handlers

import (
	"errors"
	"net/http"

	svc "github.com/PavelMilanov/stackforge/services"
	"github.com/PavelMilanov/stackforge/views/pages"
	"github.com/labstack/echo/v5"
)

/*
templateCatalog отображает каталог Portainer templates.

Входные параметры:
- c: текущий Echo context HTTP-запроса.

Возвращает:
- error: ошибка получения templates, ошибка рендера страницы или nil при успешной записи ответа.
*/
func (h *Handler) templateCatalog(c *echo.Context) error {
	templates, err := h.Svc.TemplatesList(c.Request().Context())
	if err != nil {
		return err
	}
	if len(templates) == 0 {
		return h.render(c, pages.TemplateCatalogPage(nil, pages.TemplateView{}))
	}

	return h.render(c, pages.TemplateCatalogPage(toTemplateViews(templates), toTemplateView(templates[0])))
}

/*
templatePreview возвращает HTML-фрагмент карточки выбранного шаблона.

Входные параметры:
- c: текущий Echo context HTTP-запроса; ожидает query-параметр template_id.

Возвращает:
- error: ошибка получения template, ошибка рендера фрагмента, HTTP 404 для неизвестного template_id или nil при успешной записи ответа.
*/
func (h *Handler) templatePreview(c *echo.Context) error {
	templateID := c.QueryParam("template_id")
	template, err := h.Svc.TemplateGetByID(c.Request().Context(), templateID)
	if err != nil {
		if errors.Is(err, svc.ErrTemplateNotFound) {
			return c.String(http.StatusNotFound, "template not found")
		}
		return err
	}

	return h.render(c, pages.TemplatePreview(toTemplateView(template)))
}

/*
toTemplateViews преобразует список svc.StackTemplate в список pages.TemplateView.

Входные параметры:
- items: список шаблонов из сервисного слоя.

Возвращает:
- []pages.TemplateView: список view-моделей для отображения шаблонов на странице.
*/
func toTemplateViews(items []svc.StackTemplate) []pages.TemplateView {
	views := make([]pages.TemplateView, 0, len(items))
	for _, item := range items {
		views = append(views, toTemplateView(item))
	}

	return views
}

/*
toTemplateView преобразует svc.StackTemplate в pages.TemplateView.

Входные параметры:
- item: шаблон из сервисного слоя.

Возвращает:
- pages.TemplateView: view-модель шаблона для templ-компонентов.
*/
func toTemplateView(item svc.StackTemplate) pages.TemplateView {
	services := make([]pages.ServiceView, 0, len(item.Services))
	for _, service := range item.Services {
		services = append(services, pages.ServiceView{
			Name: service.Name,
			Note: service.Note,
		})
	}

	return pages.TemplateView{
		ID:          item.ID,
		Name:        item.Title,
		Category:    item.Category,
		Description: item.Description,
		Repository:  item.Repository,
		Metadata:    item.Metadata,
		Services:    services,
	}
}
