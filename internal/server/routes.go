package server

import (
	"net/http"

	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"

	"github.com/neoxelox/zeus/internal/logger"
)

func (s *Server) addRoutes(logger *logger.Logger) error { // nolint
	logLevel := zerolog.InfoLevel
	if s.Configuration.App.Environment == Environments.DEVELOPMENT {
		logLevel = zerolog.DebugLevel
	}

	s.Instance.Pre(middleware.RemoveTrailingSlash())
	s.Instance.Use(logger.Middleware(logLevel))
	s.Instance.Use(middleware.Recover())
	s.Instance.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: s.Configuration.App.Host,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
	}))

	s.Instance.GET("/health", s.Health)

	return nil
}
