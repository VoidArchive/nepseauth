package nepse

import (
	"fmt"
	"net/http"
)

// NepseError represents different types of NEPSE API errors
type NepseError struct {
	Type    ErrorType
	Message string
	Err     error
}

// ErrorType represents the category of NEPSE error
type ErrorType string

const (
	ErrorTypeInvalidClientRequest  ErrorType = "invalid_client_request"
	ErrorTypeInvalidServerResponse ErrorType = "invalid_server_response"
	ErrorTypeTokenExpired          ErrorType = "token_expired"
	ErrorTypeNetworkError          ErrorType = "network_error"
	ErrorTypeUnauthorized          ErrorType = "unauthorized"
	ErrorTypeNotFound              ErrorType = "not_found"
	ErrorTypeRateLimit             ErrorType = "rate_limit"
	ErrorTypeInternal              ErrorType = "internal_error"
)

// Error implements the error interface
func (e *NepseError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("nepse %s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("nepse %s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *NepseError) Unwrap() error {
	return e.Err
}

// Is checks if the error is of a specific type
func (e *NepseError) Is(target error) bool {
	if target, ok := target.(*NepseError); ok {
		return e.Type == target.Type
	}
	return false
}

// NewNepseError creates a new NEPSE error
func NewNepseError(errorType ErrorType, message string, err error) *NepseError {
	return &NepseError{
		Type:    errorType,
		Message: message,
		Err:     err,
	}
}

// NewInvalidClientRequestError creates an invalid client request error
func NewInvalidClientRequestError(message string) *NepseError {
	return NewNepseError(ErrorTypeInvalidClientRequest, message, nil)
}

// NewInvalidServerResponseError creates an invalid server response error
func NewInvalidServerResponseError(message string) *NepseError {
	return NewNepseError(ErrorTypeInvalidServerResponse, message, nil)
}

// NewTokenExpiredError creates a token expired error
func NewTokenExpiredError() *NepseError {
	return NewNepseError(ErrorTypeTokenExpired, "access token expired", nil)
}

// NewNetworkError creates a network error
func NewNetworkError(err error) *NepseError {
	return NewNepseError(ErrorTypeNetworkError, "network request failed", err)
}

// NewUnauthorizedError creates an unauthorized error
func NewUnauthorizedError(message string) *NepseError {
	return NewNepseError(ErrorTypeUnauthorized, message, nil)
}

// NewNotFoundError creates a not found error
func NewNotFoundError(resource string) *NepseError {
	return NewNepseError(ErrorTypeNotFound, fmt.Sprintf("%s not found", resource), nil)
}

// NewRateLimitError creates a rate limit error
func NewRateLimitError() *NepseError {
	return NewNepseError(ErrorTypeRateLimit, "rate limit exceeded", nil)
}

// NewInternalError creates an internal error
func NewInternalError(message string, err error) *NepseError {
	return NewNepseError(ErrorTypeInternal, message, err)
}

// MapHTTPStatusToError maps HTTP status codes to NEPSE errors
func MapHTTPStatusToError(statusCode int, message string) *NepseError {
	switch statusCode {
	case http.StatusBadRequest:
		return NewInvalidClientRequestError(message)
	case http.StatusUnauthorized:
		return NewTokenExpiredError()
	case http.StatusForbidden:
		return NewUnauthorizedError(message)
	case http.StatusNotFound:
		return NewNotFoundError("resource")
	case http.StatusTooManyRequests:
		return NewRateLimitError()
	case http.StatusBadGateway:
		return NewInvalidServerResponseError(message)
	case http.StatusServiceUnavailable:
		return NewInvalidServerResponseError("service unavailable")
	case http.StatusGatewayTimeout:
		return NewInvalidServerResponseError("gateway timeout")
	default:
		if statusCode >= 500 {
			return NewInvalidServerResponseError(fmt.Sprintf("server error: %d", statusCode))
		}
		return NewInternalError(fmt.Sprintf("unexpected status code: %d", statusCode), nil)
	}
}

// IsRetryable returns true if the error is potentially retryable
func (e *NepseError) IsRetryable() bool {
	switch e.Type {
	case ErrorTypeTokenExpired, ErrorTypeNetworkError, ErrorTypeInvalidServerResponse:
		return true
	case ErrorTypeRateLimit:
		return true // With backoff
	default:
		return false
	}
}
