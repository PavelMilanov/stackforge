package handlers

import (
	"net/http"

	"github.com/PavelMilanov/stackforge/config"
	"github.com/PavelMilanov/stackforge/integrations/portainer"
	svc "github.com/PavelMilanov/stackforge/services"
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
*/
func NewHandler(env *config.Env, portainerClient *portainer.Client) *Handler {
	return &Handler{
		Env:       env,
		Templates: svc.NewPortainerService(portainerClient),
	}
}

/*
InitRouters создает Echo router, подключает middleware, статику и маршруты страниц.
*/
func (h *Handler) InitRouters() *echo.Echo {
	e := echo.New()
	//e.Use(middleware.RequestLogger())
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
registerPageRoutes описывает server-side rendered страницы приложения.
*/
func (h *Handler) registerPageRoutes(e *echo.Echo) {
	e.GET("/", h.templateCatalog)
	e.GET("/templates/preview", h.templatePreview)
	e.GET("/stands", h.stands)
	e.GET("/history", h.history)
	e.GET("/docs", h.docs)
}

/*
check возвращает минимальный healthcheck без обращения к внешним интеграциям.
*/
func (h *Handler) check(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
