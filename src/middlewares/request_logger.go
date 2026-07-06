package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v5"
	"github.com/labstack/echo/v5/middleware"
	"github.com/sirupsen/logrus"
)

/*
RequestLogger создает middleware логирования HTTP-запросов.

Входные параметры:
- отсутствуют.

Возвращает:
- echo.MiddlewareFunc: middleware Echo, который пишет logrus-запись для каждого запроса.
*/
func RequestLogger() echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		HandleError: true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogHost:     true,
		LogMethod:   true,
		LogURI:      true,
		LogURIPath:  true,
		LogStatus:   true,
		LogValuesFunc: func(_ *echo.Context, v middleware.RequestLoggerValues) error {
			fields := logrus.Fields{
				"method":    v.Method,
				"uri":       v.URI,
				"path":      v.URIPath,
				"status":    v.Status,
				"latency":   v.Latency.String(),
				"remote_ip": v.RemoteIP,
				"host":      v.Host,
			}
			if v.Status >= http.StatusInternalServerError {
				fields["error"] = http.StatusText(v.Status)
				if v.Error != nil {
					fields["error"] = v.Error.Error()
				}
			}

			entry := logrus.WithFields(fields)
			if v.Error != nil {
				entry.Error("HTTP request")
				return nil
			}
			if v.Status >= http.StatusInternalServerError {
				entry.Error("HTTP request")
				return nil
			}
			if v.Status >= http.StatusBadRequest {
				entry.Warn("HTTP request")
				return nil
			}
			entry.Info("HTTP request")
			return nil
		},
	})
}
