package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/neoxelox/zeus/pkg/payload"
	"github.com/neoxelox/zeus/pkg/user"
)

// UserHandler describes the user handler.
type UserHandler struct {
	userCreator user.CreatorUseCase
	userGetter  user.GetterUseCase
}

// NewUserHandler creates a new UserHandler instance.
func NewUserHandler(userCreator user.CreatorUseCase, userGetter user.GetterUseCase) *UserHandler {
	return &UserHandler{
		userCreator: userCreator,
		userGetter:  userGetter,
	}
}

// Create creates a new user.
func (h *UserHandler) Create(ctx echo.Context) error {
	var req payload.UserCreateRequest
	if err := ctx.Bind(&req); err != nil {
		return payload.ErrInvalidRequest.Wrap(err, "Cannot bind user create request")
	}
	if err := ctx.Validate(&req); err != nil {
		return payload.ErrInvalidRequest.Wrap(err, "Cannot validate user create request")
	}

	m, err := h.userCreator.Create(ctx.Request().Context(), req.Name, req.Username, req.Age)
	if err != nil {
		return err // nolint
	}

	res := payload.NewUserCreateResponse(m)

	return ctx.JSON(http.StatusOK, res)
}

// GetByID gets a user by its ID.
func (h *UserHandler) GetByID(ctx echo.Context) error {
	var req payload.UserGetByIDRequest
	if err := ctx.Bind(&req); err != nil {
		return payload.ErrInvalidRequest.Wrap(err, "Cannot bind user get by id request")
	}
	if err := ctx.Validate(&req); err != nil {
		return payload.ErrInvalidRequest.Wrap(err, "Cannot validate user get by id request")
	}

	m, err := h.userGetter.GetByID(ctx.Request().Context(), req.ID)
	if err != nil {
		return err // nolint
	}

	res := payload.NewUserGetByIDResponse(m)

	return ctx.JSON(http.StatusOK, res)
}

// List gets existing users with a similar username.
func (h *UserHandler) List(ctx echo.Context) error {
	var req payload.UserListRequest
	if err := ctx.Bind(&req); err != nil {
		return payload.ErrInvalidRequest.Wrap(err, "Cannot bind user list request")
	}
	if err := ctx.Validate(&req); err != nil {
		return payload.ErrInvalidRequest.Wrap(err, "Cannot validate user list request")
	}

	ms, err := h.userGetter.List(ctx.Request().Context(), req.Username)
	if err != nil {
		return err // nolint
	}

	res := payload.NewUserListResponse(ms)

	return ctx.JSON(http.StatusOK, res)
}
