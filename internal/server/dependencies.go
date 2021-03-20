package server

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"

	"github.com/neoxelox/zeus/internal/database"
	"github.com/neoxelox/zeus/internal/logger"
)

// Dependencies describes the application dependencies.
type Dependencies struct {
	Database *database.Database
}

func (s *Server) addDependencies(logger *logger.Logger) error {
	zlogLevel := zerolog.InfoLevel
	plogLevel := pgx.LogLevel(pgx.LogLevelError)
	if s.Configuration.App.Environment == Environments.DEVELOPMENT {
		zlogLevel = zerolog.DebugLevel
		plogLevel = pgx.LogLevelDebug
	}

	database, err := database.New(context.Background(), s.Configuration.Database.Dsn, 15, database.Configuration{
		MinConns: 0,
		MaxConns: 22, // nolint
		AppName:  s.Configuration.App.Name,
		Logger:   logger.Database(zlogLevel),
		LogLevel: plogLevel,
	})
	if err != nil {
		return errors.Wrap(err, "Cannot add database dependency")
	}

	s.Dependencies = Dependencies{
		Database: database,
	}

	return nil
}
