package user

import (
	"context"

	"github.com/cockroachdb/errors"

	"github.com/neoxelox/zeus/internal/database"
	"github.com/neoxelox/zeus/pkg/model"
	"github.com/neoxelox/zeus/pkg/repository"
)

// CreatorUseCase interacts with the user creator use case.
type CreatorUseCase interface {
	Create(ctx context.Context, name string, username string, age int) (*model.User, error)
}

// Creator implements the CreatorUseCase.
type Creator struct {
	userRepository repository.UserRepository
}

// NewCreator creates a new Creator instance.
func NewCreator(userRepository repository.UserRepository) *Creator {
	return &Creator{
		userRepository: userRepository,
	}
}

// Create creates a new user.
func (c *Creator) Create(ctx context.Context, name string, username string, age int) (*model.User, error) {
	if age < model.UserMinAge {
		return nil, model.ErrUserBelowAge.New("Cannot create user underaged")
	}

	user := model.NewUser(name, username, age)

	user, err := c.userRepository.Create(ctx, user)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrIntegrityViolation):
			return nil, model.ErrExistingUsername.Wrap(err, "Cannot create user with existing username")
		default:
			return nil, errors.Wrap(err, "Cannot create user")
		}
	}

	return user, nil
}
