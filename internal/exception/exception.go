package exception

import (
	"fmt"
	"math/rand"
	"net/http"
	"runtime"

	"github.com/cockroachdb/errors"
	"github.com/labstack/echo/v4"
)

type exception interface {
	origin() int
	inner() error
	status() int
	message() string
	Wrap(err error, msg string) error
	New(msg string) error
	Is(reference error) bool
	Error() string
	Unwrap() error
	String() string
}

// Exception describes a complex error.
type Exception struct {
	_origin int    `json:"-"`
	_inner  error  `json:"-"`
	Status  int    `json:"-"`
	Message string `json:"message"`
}

// New creates a new Exception instance.
func New(status int, message string) Exception {
	origin := rand.Intn(2048) // nolint
	if pc, _, _, ok := runtime.Caller(1); ok {
		origin = int(pc)
	}

	return Exception{
		_origin: origin,
		Status:  status,
		Message: message,
	}
}

func (e Exception) origin() int {
	return e._origin
}

func (e Exception) inner() error {
	return e._inner
}

func (e Exception) status() int {
	return e.Status
}

func (e Exception) message() string {
	return e.Message
}

// Wrap wraps err with msg in exception.
func (e Exception) Wrap(err error, msg string) error {
	return Exception{
		_origin: e._origin,
		_inner:  errors.Wrap(err, msg),
		Status:  e.Status,
		Message: e.Message,
	}
}

// New wraps msg in exception.
func (e Exception) New(msg string) error {
	return Exception{
		_origin: e._origin,
		_inner:  errors.New(msg),
		Status:  e.Status,
		Message: e.Message,
	}
}

// Is checks whether the exception is the given reference error.
func (e Exception) Is(reference error) bool {
	if other, ok := reference.(exception); ok { // nolint
		return e._origin == other.origin()
	}

	return false
}

// Error satisfies the standard error interface.
func (e Exception) Error() string {
	return e.Message
}

// Unwrap satisfies the standard error interface.
func (e Exception) Unwrap() error {
	return e._inner
}

// String satisfies the standard error interface.
func (e Exception) String() string {
	return fmt.Sprintf("<%s: %d>", e.Message, e.Status)
}

// ErrGeneric generic error.
var ErrGeneric = New(http.StatusInternalServerError, "ERR_GENERIC")

// Handler controls all the returned exceptions.
func Handler(err error, ctx echo.Context) {
	if exc, ok := err.(exception); ok { // nolint
		ctx.Logger().Error(exc.inner()) // TODO(alex): Send to Sentry with ctx.
		returnException(exc, ctx)
	} else { // Fallback.
		ctx.Logger().Error(err)
		returnException(ErrGeneric, ctx)
	}
}

func returnException(exc exception, ctx echo.Context) {
	var err error

	if ctx.Response().Committed {
		return
	}

	if ctx.Request().Method == http.MethodHead {
		err = ctx.NoContent(exc.status())
	} else {
		err = ctx.JSON(exc.status(), exc)
	}

	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, "Cannot return exception"))
	}
}
