package handlers

import (
	"errors"
	"net/http"

	svc "github.com/PavelMilanov/stackforge/services"
	"github.com/PavelMilanov/stackforge/views/pages"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/sirupsen/logrus"
)

/*
render пишет HTTP 200 и рендерит templ-компонент в response writer Echo.
*/
func (h *Handler) render(c *echo.Context, component templ.Component) error {
	c.Response().WriteHeader(http.StatusOK)
	return component.Render(c.Request().Context(), c.Response())
}

/*
templateCatalog отображает каталог Portainer templates.
*/
func (h *Handler) templateCatalog(c *echo.Context) error {
	templates, err := h.Templates.List(c.Request().Context())
	if err != nil {
		logrus.WithError(err).Error("Ошибка")
		return err
	}
	if len(templates) == 0 {
		return h.render(c, pages.TemplateCatalogPage(nil, pages.TemplateView{}))
	}

	return h.render(c, pages.TemplateCatalogPage(toTemplateViews(templates), toTemplateView(templates[0])))
}

/*
templatePreview возвращает HTML-фрагмент карточки выбранного шаблона.
*/
func (h *Handler) templatePreview(c *echo.Context) error {
	templateID := c.QueryParam("template_id")
	template, err := h.Templates.GetByID(c.Request().Context(), templateID)
	if err != nil {
		if errors.Is(err, svc.ErrTemplateNotFound) {
			return c.String(http.StatusNotFound, "template not found")
		}
		logrus.WithError(err).Error("Ошибка")
		return err
	}

	return h.render(c, pages.TemplatePreview(toTemplateView(template)))
}

// stands отображает 4 рабочих стенда и развернутые в них стеки.
func (h *Handler) stands(c *echo.Context) error {
	return h.render(c, pages.StandsPage())
}

// history отображает журнал запусков и операций со стендами.
func (h *Handler) history(c *echo.Context) error {
	return h.render(c, pages.HistoryPage())
}

// docs отображает короткие карточки документации со ссылками на внутреннюю wiki.
func (h *Handler) docs(c *echo.Context) error {
	return h.render(c, pages.DocsPage())
}

/*
toTemplateViews преобразует список svc.StackTemplate в список pages.TemplateView.
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
		Name:        item.Name,
		Category:    item.Category,
		Description: item.Description,
		Purpose:     item.Purpose,
		Fit:         item.Fit,
		Parameters:  item.Parameters,
		Services:    services,
	}
}
