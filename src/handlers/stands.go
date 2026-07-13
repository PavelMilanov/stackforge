package handlers

import (
	"errors"
	"net/http"
	"time"

	svc "github.com/PavelMilanov/stackforge/services"
	"github.com/PavelMilanov/stackforge/views/pages"
	"github.com/labstack/echo/v5"
)

/*
stands отображает 4 рабочих стенда и развернутые в них стеки.

Входные параметры:
- c: текущий Echo context HTTP-запроса.

Возвращает:
- error: ошибка рендера страницы или nil при успешной записи ответа.
*/
func (h *Handler) stands(c *echo.Context) error {
	return h.render(c, pages.StandsPage())
}

/*
createStackModal возвращает HTML-фрагмент формы создания стека в выбранном стенде.

Входные параметры:
- c: текущий Echo context HTTP-запроса; ожидает query-параметр stand.

Возвращает:
- error: ошибка получения templates, ошибка рендера модального окна или nil при успешной записи ответа.
*/
func (h *Handler) createStackModal(c *echo.Context) error {
	templates, err := h.Svc.TemplatesList(c.Request().Context())
	if err != nil {
		return err
	}

	return h.render(c, pages.CreateStackModal(c.QueryParam("stand"), toTemplateViews(templates)))
}

/*
createStand создает стенд целиком.

Сейчас handler содержит заглушку: состав стенда и интеграционную логику создания
стеков в Portainer/Gitea нужно добавить на месте TODO.

Входные параметры:
  - c: текущий Echo context POST /stands/create; ожидает query-параметр stand.

Возвращает:
- error: ошибка обновления данных стенда, ошибка рендера результата или nil при успешной записи ответа.
*/
func (h *Handler) createStand(c *echo.Context) error {
	standNumber := c.QueryParam("stand")

	// TODO: здесь будет создание всех стеков стенда в Portainer/Gitea.
	stand, err := h.standViewFromPortainer(c, standNumber)
	if err != nil {
		return err
	}
	return h.render(c, pages.StandCard(stand))
}

/*
createStack создает один стек в выбранном стенде.

Сейчас handler содержит заглушку: реальную интеграционную логику создания стека
в Portainer/Gitea нужно добавить на месте TODO.

Входные параметры:
  - c: текущий Echo context POST /stands/stacks/create; ожидает form-параметры stand и template_id,
    которые HTMX отправляет из формы create-stack-form и handler читает через c.FormValue.

Возвращает:
- error: ошибка получения template, HTTP 404 для неизвестного template_id, ошибка обновления данных стенда, ошибка рендера результата или nil при успешной записи ответа.
*/
func (h *Handler) createStack(c *echo.Context) error {
	standNumber := c.FormValue("stand")
	templateID := c.FormValue("template_id")

	template, err := h.Svc.TemplateGetByID(c.Request().Context(), templateID)
	if err != nil {
		if errors.Is(err, svc.ErrTemplateNotFound) {
			return c.String(http.StatusNotFound, "template not found")
		}
		return err
	}

	_ = template
	// TODO: здесь будет создание стека в Portainer/Gitea.
	// template содержит выбранный шаблон, repository и metadata для создания стека.
	stand, err := h.standViewFromPortainer(c, standNumber)
	if err != nil {
		return err
	}
	return h.render(c, pages.CreateStackResult(stand))
}

/*
standStacks возвращает фрагмент списка стеков стенда, загруженный из Portainer API.

Входные параметры:
- c: текущий Echo context HTTP-запроса; ожидает query-параметр stand.

Возвращает:
- error: ошибка получения стеков из Portainer, ошибка рендера фрагмента или nil при успешной записи ответа.
*/
func (h *Handler) standStacks(c *echo.Context) error {
	standNumber := c.QueryParam("stand")
	stacks, err := h.Svc.StacksByStand(c.Request().Context(), standNumber)
	if err != nil {
		return err
	}
	return h.render(c, pages.StandStacksResult(standNumber, stacks, standCreateDate(stacks)))
}

/*
closeStandModal очищает контейнер HTMX-модалки.

Входные параметры:
- c: текущий Echo context HTTP-запроса.

Возвращает:
- error: ошибка записи строкового ответа Echo или nil при успешной очистке контейнера.
*/
func (h *Handler) closeStandModal(c *echo.Context) error {
	return c.String(http.StatusOK, "")
}

/*
standViewFromPortainer собирает view-модель стенда по данным Portainer.

Входные параметры:
- c: текущий Echo context HTTP-запроса.
- standNumber: номер стенда, по которому фильтруются стеки.

Возвращает:
- pages.StandView: view-модель стенда со списком стеков и датой создания.
- error: ошибка получения стеков из Portainer или nil при успешной сборке модели.
*/
func (h *Handler) standViewFromPortainer(c *echo.Context, standNumber string) (pages.StandView, error) {
	stacks, err := h.Svc.StacksByStand(c.Request().Context(), standNumber)
	if err != nil {
		return pages.StandView{}, err
	}

	return pages.StandWithStacks(standNumber, stacks, standCreateDate(stacks)), nil
}

/*
standCreateDate вычисляет дату создания стенда по самому раннему времени создания его стеков.

Входные параметры:
- stacks: список стеков одного стенда.

Возвращает:
- string: строка для отображения даты создания стенда или "Дата создания: -", если дату определить нельзя.
*/
func standCreateDate(stacks []svc.StackInfo) string {
	if len(stacks) == 0 {
		return "Дата создания: -"
	}

	createdAt := stacks[0].CreatedAt
	for _, stack := range stacks[1:] {
		if stack.CreatedAt.Before(createdAt) {
			createdAt = stack.CreatedAt
		}
	}
	if createdAt.IsZero() || createdAt.Equal(time.Unix(0, 0)) {
		return "Дата создания: -"
	}

	return "Дата создания: " + createdAt.Format("02.01.2006 15:04")
}
