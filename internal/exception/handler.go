package exception

import (
	"github.com/labstack/echo/v4"
)

// Handler controls all the returned exceptions.
func Handler(err error, ctx echo.Context) {
	// TODO(alex): create custom global error handler and app exceptions.

	ctx.Echo().DefaultHTTPErrorHandler(err, ctx)
}
