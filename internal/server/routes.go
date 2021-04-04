package server

import (
	"fmt"
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

	allowOrigins := []string{}
	for _, origin := range s.Configuration.App.Host {
		allowOrigins = append(allowOrigins, fmt.Sprintf("%s://%s", s.Configuration.App.Scheme, origin))
	}

	s.Instance.Pre(middleware.RemoveTrailingSlash()) // TODO(alex): Move to Horae.
	s.Instance.Use(logger.Middleware(logLevel))
	s.Instance.Use(middleware.Recover())
	s.Instance.Use(middleware.CORSWithConfig(middleware.CORSConfig{ // TODO(alex): Move to Horae.
		AllowOrigins: allowOrigins,
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut},
		AllowHeaders: []string{"*"},
		MaxAge:       86400, // nolint
	}))
	s.Instance.Use(middleware.BodyLimit("2M")) // TODO(alex): Move to Horae.

	// Endpoints.

	s.Instance.GET("/health", s.Health)

	v1 := s.Instance.Group("/v1")

	user := v1.Group("/user")
	user.GET("", s.Handlers.User.List)
	user.GET("/:id", s.Handlers.User.GetByID)
	user.POST("", s.Handlers.User.Create)

	return nil
}
