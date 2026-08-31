package health

import (
	"context"
	"net/http"
	"time"

	"github.com/labstack/echo/v5"
)

const healthCheckTimeout = time.Second

type databasePinger interface {
	Ping(context.Context) error
}

type healthResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func Handler(database databasePinger) echo.HandlerFunc {
	return func(c *echo.Context) error {
		ctx, cancel := context.WithTimeout(c.Request().Context(), healthCheckTimeout)
		defer cancel()

		if err := database.Ping(ctx); err != nil {
			return c.JSON(http.StatusServiceUnavailable, errorResponse{
				Message: "database unavailable",
				Code:    "service_unavailable",
			})
		}

		return c.JSON(http.StatusOK, healthResponse{Status: "ok"})
	}
}
