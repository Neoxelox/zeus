package server

import (
	"context"
	"fmt"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"

	"github.com/neoxelox/zeus/internal/exception"
	"github.com/neoxelox/zeus/internal/logger"
)

// Server describes the main application instance.
type Server struct {
	Instance      *echo.Echo
	Configuration Configuration
	Dependencies  Dependencies
	Handlers      Handlers
}

// NewServer creates a new Server instance.
func New(e *echo.Echo) *Server {
	server := &Server{Instance: e}

	if err := server.addConfiguration(); err != nil {
		panic("Cannot add server configuration")
	}

	debug := false
	logLevel := zerolog.InfoLevel
	if server.Configuration.App.Environment == Environments.DEVELOPMENT {
		debug = true
		logLevel = zerolog.DebugLevel
	}

	appLogger := logger.New(server.Configuration.App.Name)
	server.Instance.Logger = appLogger.Standard(logLevel)
	server.Instance.HideBanner = true
	server.Instance.HidePort = true
	server.Instance.Debug = debug
	server.Instance.HTTPErrorHandler = exception.Handler
	// Server.Instance.Renderer
	// TODO(alex): add custom server.Instance.Validator.
	server.Instance.IPExtractor = echo.ExtractIPFromRealIPHeader()

	if err := server.addDependencies(appLogger); err != nil {
		panic("Cannot add server dependencies")
	}

	if err := server.addHandlers(); err != nil {
		panic("Cannot add server handlers")
	}

	if err := server.addRoutes(appLogger); err != nil {
		panic("Cannot add server routes")
	}

	return server
}

// Startup starts the server.
func (s *Server) Startup() {
	s.Instance.Logger.Info("Server startup")
	s.Instance.Logger.Fatal(s.Instance.Start(fmt.Sprintf(":%d", s.Configuration.App.Port)))
}

// Shutdown stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	// Deadline := time.Duration(s.Configuration.App.GracefulTimeout) * time.Second
	// if ctxDeadline, ok := ctx.Deadline(); ok {
	// 	deadline = time.Until(ctxDeadline)
	// }.

	if err := s.Dependencies.Database.Close(ctx); err != nil {
		return errors.Wrap(err, "Cannot close connection to the database")
	}

	s.Instance.Logger.Info("Server shutdown")

	appLogger, _ := s.Instance.Logger.(*logger.Logger)
	if err := appLogger.Flush(); err != nil {
		return errors.Wrap(err, "Cannot flush main logger instance")
	}

	if err := s.Instance.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "Cannot stop main application instance")
	}

	return nil
}
