package model

import (
	"net/http"
	"time"

	"github.com/rs/xid"

	"github.com/neoxelox/zeus/internal/exception"
)

// User represents a user.
type User struct {
	ID        xid.ID     `json:"id" db:"id"`
	Name      string     `json:"name" db:"name"`
	Username  string     `json:"username" db:"username"`
	Age       int        `json:"age" db:"age"`
	CreatedAt time.Time  `json:"-" db:"created_at"`
	UpdatedAt time.Time  `json:"-" db:"updated_at"`
	DeletedAt *time.Time `json:"-" db:"deleted_at"`
}

// NewUser creates a new User instance.
func NewUser(name string, username string, age int) *User {
	now := time.Now()

	return &User{
		ID:        xid.New(),
		Name:      name,
		Username:  username,
		Age:       age,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// UserMinAge minimum age for user to exist.
const UserMinAge = 18

var (
	// ErrUserBelowAge user is below UserMinAge.
	ErrUserBelowAge = exception.New(http.StatusBadRequest, "ERR_USER_BELOW_AGE")

	// ErrExistingUsername username already exists.
	ErrExistingUsername = exception.New(http.StatusBadRequest, "ERR_EXISTING_USERNAME")

	// ErrUserNotExists user not exists.
	ErrUserNotExists = exception.New(http.StatusBadRequest, "ERR_USER_NOT_EXISTS")
)
