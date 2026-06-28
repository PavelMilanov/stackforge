package handlers

import (
	"context"
	"net/http"

	"github.com/PavelMilanov/stackforge/config"
	svc "github.com/PavelMilanov/stackforge/services"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

// Handler объединяет конфигурацию приложения и HTTP-обработчики.
// Позже сюда будут добавлены сервисы интеграций Portainer, Gitea и storage-слой.
type Handler struct {
	Env       *config.Env
	Templates TemplateService
}

// TemplateService описывает источник Portainer templates для HTTP-слоя.
// Сейчас используется static-реализация, позже ее заменит клиент Portainer API.
type TemplateService interface {
	List(ctx context.Context) ([]svc.StackTemplate, error)
	GetByID(ctx context.Context, id string) (svc.StackTemplate, error)
}

// NewHandler создает HTTP handler с runtime-конфигурацией приложения.
func NewHandler(env *config.Env) *Handler {
	return &Handler{
		Env:       env,
		Templates: svc.NewStaticService(),
	}
}

// InitRouters создает Echo router, подключает middleware, статику и маршруты страниц.
func (h *Handler) InitRouters() *echo.Echo {
	e := echo.New()

	// CORS сейчас открыт для локального этапа разработки интерфейса.
	// Перед production-запуском список origins должен быть ограничен доменом StackForge.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Static assets собираются frontend pipeline и отдаются без CDN.
	e.Static("/assets", "public/assets")

	// /check используется Docker healthcheck и reverse proxy для проверки доступности сервиса.
	e.GET("/check", h.check)
	h.registerPageRoutes(e)

	return e
}

// registerPageRoutes описывает server-side rendered страницы приложения.
// Маршруты держатся отдельно от InitRouters, чтобы HTTP setup и web-страницы не смешивались.
func (h *Handler) registerPageRoutes(e *echo.Echo) {
	e.GET("/", h.templateCatalog)
	e.GET("/templates/preview", h.templatePreview)
	e.GET("/stands", h.stands)
	e.GET("/history", h.history)
	e.GET("/docs", h.docs)
}

// check возвращает минимальный healthcheck без обращения к внешним интеграциям.
func (h *Handler) check(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
