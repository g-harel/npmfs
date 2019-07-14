package registry

import (
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

// Error represents a generic registry error.
// It is used to determine if an error xerrors.Is() a registry error.
var Error = xerrors.New("registry error")

// Err indicates the appropriate status code return to the user.
type Err struct {
	StatusCode int
}

var _ error = &Err{}

// Is returns true if the given err is the generic registry.Error.
// It implements the interface described in xerrors.
func (e *Err) Is(err error) bool {
	return err == Error
}

// As returns true if the given err is the generic registry.Error.
// It implements the interface described in xerrors.
func (e *Err) As(err error, target interface{}) bool {
	if err == Error {
		target = e
		return true
	}
	return false
}

// Error implements the error interface and formats the
func (e *Err) Error() string {
	return fmt.Sprintf("%v %v", e.StatusCode, strings.ToLower(http.StatusText(e.StatusCode)))
}

// ErrNotFound signals a package and version combination is not found.
var ErrNotFound = &Err{http.StatusNotFound}

// ErrGatewayTimeout signals the remote registry did not respond before the timeout duration.
var ErrGatewayTimeout = &Err{http.StatusGatewayTimeout}

// ErrBadGateway signals an unexpected upstream error.
var ErrBadGateway = &Err{http.StatusBadGateway}
