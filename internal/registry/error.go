package registry

import (
	"net/http"
)

// Error indicates the appropriate status code return to the user.
type Error struct {
	StatusCode int
}

// Error implements the error interface and provides a short description of the error.
func (e *Error) Error() string {
	return http.StatusText(e.StatusCode)
}

// ErrNotFound signals a package and version combination is not found.
var ErrNotFound = &Error{http.StatusNotFound}

// ErrGatewayTimeout signals the remote registry did not respond before the timeout duration.
var ErrGatewayTimeout = &Error{http.StatusGatewayTimeout}

// ErrBadGateway signals an unexpected upstream error.
var ErrBadGateway = &Error{http.StatusBadGateway}
