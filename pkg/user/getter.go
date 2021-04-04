package user

import (
	"context"

	"github.com/cockroachdb/errors"
	"github.com/rs/xid"

	"github.com/neoxelox/zeus/internal/database"
	"github.com/neoxelox/zeus/pkg/model"
	"github.com/neoxelox/zeus/pkg/repository"
)

// GetterUseCase interacts with the user getter use case.
type GetterUseCase interface {
	GetByID(ctx context.Context, ID xid.ID) (*model.User, error)
	List(ctx context.Context, username string) ([]model.User, error)
}

// Getter implements the GetterUseCase.
type Getter struct {
	userRepository repository.UserRepository
}

// NewGetter creates a new Getter instance.
func NewGetter(userRepository repository.UserRepository) *Getter {
	return &Getter{
		userRepository: userRepository,
	}
}

// GetByID gets a user by its ID.
func (g *Getter) GetByID(ctx context.Context, ID xid.ID) (*model.User, error) {
	user, err := g.userRepository.GetByID(ctx, ID)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNoRows):
			return nil, model.ErrUserNotExists.Wrap(err, "Cannot get a user with that id")
		default:
			return nil, errors.Wrap(err, "Cannot get user")
		}
	}

	return user, nil
}

// List gets existing users with a similar username.
func (g *Getter) List(ctx context.Context, username string) ([]model.User, error) {
	users, err := g.userRepository.List(ctx, username)
	if err != nil {
		return nil, errors.Wrap(err, "Cannot list users")
	}

	return users, nil
}
