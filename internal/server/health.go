package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) Health(ctx echo.Context) error {
	err := s.Dependencies.Database.Health(ctx.Request().Context())
	if err != nil {
		ctx.Logger().Error(err)

		return echo.ErrServiceUnavailable
	}

	return ctx.String(http.StatusOK, "OK\n")
}
