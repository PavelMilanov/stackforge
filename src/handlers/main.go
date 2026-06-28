package handlers

import (
	"net/http"

	"github.com/PavelMilanov/stackforge/config"
	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
)

type Handler struct {
	Env *config.Env
}

/*
NewHandler создает обработчик HTTP-запросов с конфигурацией и storage-клиентом.

Параметры:
  - env: конфигурация приложения.
  - s3: storage-клиент для чтения объектов.

Результат: экземпляр Handler.
*/
func NewHandler(env *config.Env) *Handler {
	return &Handler{Env: env}
}

/*
InitRouters инициализирует Echo router и регистрирует маршруты сервиса.

Параметры: текущий экземпляр Handler.
Результат: настроенный Echo router.
*/
func (h *Handler) InitRouters() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Static("/assets", "public/assets")
	e.GET("/check", h.check)
	h.registerPageRoutes(e)

	return e
}

func (h *Handler) check(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]any{"status": "ok"})
}
