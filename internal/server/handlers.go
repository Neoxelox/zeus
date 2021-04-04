package server

import (
	"github.com/neoxelox/zeus/pkg/handler"
	"github.com/neoxelox/zeus/pkg/repository"
	"github.com/neoxelox/zeus/pkg/user"
)

// Handlers describes the application handlers.
type Handlers struct {
	User handler.UserHandler
}

func (s *Server) addHandlers() error { // nolint
	// Repositories.

	userDatabase := repository.NewUserDatabase(s.Dependencies.Database.Connection)

	// Use Cases.

	userCreator := user.NewCreator(userDatabase)
	userGetter := user.NewGetter(userDatabase)

	// Handlers.

	userHandler := handler.NewUserHandler(userCreator, userGetter)

	// Add to server.

	s.Handlers = Handlers{
		User: *userHandler,
	}

	return nil
}
