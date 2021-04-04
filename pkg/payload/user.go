package payload

import (
	"github.com/rs/xid"

	"github.com/neoxelox/zeus/pkg/model"
)

type (
	// UserCreateRequest describes the user create request.
	UserCreateRequest struct {
		Name     string `json:"name" validate:"required"`
		Username string `json:"username" validate:"required"`
		Age      int    `json:"age" validate:"required"`
	}

	// UserCreateResponse describes the user create response.
	UserCreateResponse struct {
		User model.User `json:"user"`
	}
)

// NewUserCreateResponse creates a new UserCreateResponse instance.
func NewUserCreateResponse(m *model.User) *UserCreateResponse {
	return &UserCreateResponse{
		User: *m,
	}
}

type (
	// UserGetByIDRequest describes the user get by id request.
	UserGetByIDRequest struct {
		ID xid.ID `param:"id" validate:"required"`
	}

	// UserGetByIDResponse describes the user get by id response.
	UserGetByIDResponse struct {
		User model.User `json:"user"`
	}
)

// NewUserCreateResponse creates a new UserGetByIDResponse instance.
func NewUserGetByIDResponse(m *model.User) *UserGetByIDResponse {
	return &UserGetByIDResponse{
		User: *m,
	}
}

type (
	// UserListRequest describes the user list request.
	UserListRequest struct {
		Username string `query:"username" validate:"required"`
	}

	// UserListResponse describes the user list response.
	UserListResponse struct {
		Users []model.User `json:"users"`
	}
)

// NewUserListResponse creates a new UserListResponse instance.
func NewUserListResponse(ms []model.User) *UserListResponse {
	if len(ms) == 0 {
		ms = make([]model.User, 0)
	}

	return &UserListResponse{
		Users: ms,
	}
}
