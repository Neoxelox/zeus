package payload

import (
	"net/http"

	"github.com/neoxelox/zeus/internal/exception"
)

// ErrInvalidRequest invalid headers, parameters or body for request.
var ErrInvalidRequest = exception.New(http.StatusBadRequest, "ERR_INVALID_REQUEST")
