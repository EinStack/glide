package schemas

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

type ErrorName = string

var (
	// General API errors
	UnsupportedMediaType ErrorName = "http.unsupported_media_type"
	RouteNotFound        ErrorName = "http.not_found"
	PayloadParseError    ErrorName = "http.payload_parse_error"
	UnknownError         ErrorName = "http.unknown_error"

	// Router-specific errors
	RouterNotFound       ErrorName = "routers.not_found"
	NoModelConfigured    ErrorName = "routers.no_model_configured"
	ModelUnavailable     ErrorName = "routers.model_unavailable"
	AllModelsUnavailable ErrorName = "routers.all_models_unavailable"
)

// Error / Error contains more context than the built-in error type,
// so we know information like error code and message that are useful to propagate to clients
type Error struct {
	Status  int    `json:"-"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

var _ error = &Error{}

// Error returns the error message.
func (e Error) Error() string {
	return fmt.Sprintf("Error (%s): %s", e.Name, e.Message)
}

func NewError(status int, name string, message string) Error {
	return Error{Status: status, Name: name, Message: message}
}

var ErrUnsupportedMediaType = NewError(
	fiber.StatusBadRequest,
	UnsupportedMediaType,
	"application/json is the only supported media type",
)

var ErrRouteNotFound = NewError(
	fiber.StatusNotFound,
	RouteNotFound,
	"requested route is not found",
)

var ErrRouterNotFound = NewError(fiber.StatusNotFound, RouterNotFound, "router is not found")

var ErrNoModelAvailable = NewError(
	503,
	AllModelsUnavailable,
	"all providers are unavailable",
)

func NewPayloadParseErr(err error) Error {
	return NewError(
		fiber.StatusBadRequest,
		PayloadParseError,
		err.Error(),
	)
}

func FromErr(err error) Error {
	if apiErr, ok := err.(Error); ok {
		return apiErr
	}

	return NewError(
		fiber.StatusInternalServerError,
		UnknownError,
		err.Error(),
	)
}
