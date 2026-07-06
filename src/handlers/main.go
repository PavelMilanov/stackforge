package handlers

import (
	"net/http"

	"github.com/PavelMilanov/stackforge/config"
	"github.com/PavelMilanov/stackforge/integrations/portainer"
	appmw "github.com/PavelMilanov/stackforge/middlewares"
	svc "github.com/PavelMilanov/stackforge/services"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

/*
Handler объединяет конфигурацию приложения и HTTP-обработчики.
*/
type Handler struct {
	Env       *config.Env
	Templates *svc.PortainerService
}

/*
NewHandler создает HTTP handler с runtime-конфигурацией приложения.

Входные параметры:
- env: конфигурация окружения приложения.
- portainerClient: клиент для работы с Portainer API.

Возвращает:
- *Handler: handler с подключенным сервисом Portainer templates.
*/
func NewHandler(env *config.Env, portainerClient *portainer.Client) *Handler {
	return &Handler{
		Env:       env,
		Templates: svc.NewPortainerService(portainerClient),
	}
}

/*
InitRouters создает Echo router, подключает middleware, статику и маршруты страниц.

Входные параметры:
- отсутствуют.

Возвращает:
- *echo.Echo: настроенный HTTP router приложения.
*/
func (h *Handler) InitRouters() *echo.Echo {
	e := echo.New()
	e.Use(appmw.RequestLogger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{h.Env.Cors.Origin},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Static assets собираются frontend pipeline и отдаются без CDN.
	e.Static("/assets", "public/assets")

	// /check используется Docker healthcheck и reverse proxy для проверки доступности сервиса.
	e.GET("/check", h.check)
	h.registerPageRoutes(e)

	return e
}

/*
render пишет HTTP 200 и рендерит templ-компонент в response writer Echo.

Входные параметры:
- c: текущий Echo context HTTP-запроса.
- component: templ-компонент, который нужно отрисовать в ответ.

Возвращает:
- error: ошибка рендера компонента или nil при успешной записи ответа.
*/
func (h *Handler) render(c *echo.Context, component templ.Component) error {
	c.Response().WriteHeader(http.StatusOK)
	return component.Render(c.Request().Context(), c.Response())
}
