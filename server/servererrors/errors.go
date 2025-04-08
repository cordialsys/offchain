package servererrors

import (
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// ErrorResponse represents a standardized error response
type ErrorResponse struct {
	Code       int    `json:"code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	httpStatus int    `json:"-"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf("code: %d, status: %s, message: %s", e.Code, e.Status, e.Message)
}
func (e *ErrorResponse) Send(c *fiber.Ctx) error {
	c.Status(e.httpStatus)
	return c.JSON(e)
}

// HTTP status code to gRPC code mapping
// Based on https://github.com/grpc/grpc/blob/master/doc/http-grpc-status-mapping.md
const (
	// OK is returned on success.
	CodeOK = 0
	// Canceled indicates the operation was canceled.
	CodeCanceled = 1
	// Unknown error.
	CodeUnknown = 2
	// InvalidArgument indicates client specified an invalid argument.
	CodeInvalidArgument = 3
	// DeadlineExceeded means operation expired before completion.
	CodeDeadlineExceeded = 4
	// NotFound means some requested entity was not found.
	CodeNotFound = 5
	// AlreadyExists means an attempt to create an entity failed because one already exists.
	CodeAlreadyExists = 6
	// PermissionDenied indicates the caller does not have permission to execute the specified operation.
	CodePermissionDenied = 7
	// ResourceExhausted indicates some resource has been exhausted.
	CodeResourceExhausted = 8
	// FailedPrecondition indicates operation was rejected because the system is not in a state required for the operation's execution.
	CodeFailedPrecondition = 9
	// Aborted indicates the operation was aborted.
	CodeAborted = 10
	// OutOfRange means operation was attempted past the valid range.
	CodeOutOfRange = 11
	// Unimplemented indicates operation is not implemented or not supported/enabled.
	CodeUnimplemented = 12
	// Internal errors.
	CodeInternal = 13
	// Unavailable indicates the service is currently unavailable.
	CodeUnavailable = 14
	// DataLoss indicates unrecoverable data loss or corruption.
	CodeDataLoss = 15
	// Unauthenticated indicates the request does not have valid authentication credentials.
	CodeUnauthenticated = 16
)

// Map HTTP status codes to gRPC status codes
func httpToGRPCCode(httpStatus int) int {
	switch httpStatus {
	case http.StatusOK:
		return CodeOK
	case http.StatusBadRequest:
		return CodeInvalidArgument
	case http.StatusUnauthorized:
		return CodeUnauthenticated
	case http.StatusForbidden:
		return CodePermissionDenied
	case http.StatusNotFound:
		return CodeNotFound
	case http.StatusConflict:
		return CodeAlreadyExists
	case http.StatusTooManyRequests:
		return CodeResourceExhausted
	case http.StatusNotImplemented:
		return CodeUnimplemented
	case http.StatusServiceUnavailable:
		return CodeUnavailable
	case http.StatusGatewayTimeout:
		return CodeDeadlineExceeded
	default:
		if httpStatus >= 400 && httpStatus < 500 {
			return CodeInvalidArgument
		}
		return CodeInternal
	}
}

// Map gRPC status codes to string representations
func statusCodeToString(code int) string {
	switch code {
	case CodeOK:
		return "OK"
	case CodeCanceled:
		return "Canceled"
	case CodeUnknown:
		return "Unknown"
	case CodeInvalidArgument:
		return "InvalidArgument"
	case CodeDeadlineExceeded:
		return "DeadlineExceeded"
	case CodeNotFound:
		return "NotFound"
	case CodeAlreadyExists:
		return "AlreadyExists"
	case CodePermissionDenied:
		return "PermissionDenied"
	case CodeResourceExhausted:
		return "ResourceExhausted"
	case CodeFailedPrecondition:
		return "FailedPrecondition"
	case CodeAborted:
		return "Aborted"
	case CodeOutOfRange:
		return "OutOfRange"
	case CodeUnimplemented:
		return "Unimplemented"
	case CodeInternal:
		return "Internal"
	case CodeUnavailable:
		return "Unavailable"
	case CodeDataLoss:
		return "DataLoss"
	case CodeUnauthenticated:
		return "Unauthenticated"
	default:
		return "Unknown"
	}
}

// sendError sends a standardized error response
func sendError(httpStatus int, message string) error {
	code := httpToGRPCCode(httpStatus)
	status := statusCodeToString(code)

	return &ErrorResponse{
		Code:       code,
		Status:     status,
		Message:    message,
		httpStatus: httpStatus,
	}
}

func NewErrorf(code int, msg string, args ...interface{}) error {
	return sendError(code, fmt.Sprintf(msg, args...))
}

// BadRequestf sends a 400 Bad Request error with formatted message
func BadRequestf(format string, args ...interface{}) error {
	return NewErrorf(http.StatusBadRequest, format, args...)
}

// Unauthorizedf sends a 401 Unauthorized error with formatted message
func Unauthorizedf(format string, args ...interface{}) error {
	return NewErrorf(http.StatusUnauthorized, format, args...)
}

// Forbiddenf sends a 403 Forbidden error with formatted message
func Forbiddenf(format string, args ...interface{}) error {
	return NewErrorf(http.StatusForbidden, format, args...)
}

// NotFoundf sends a 404 Not Found error with formatted message
func NotFoundf(format string, args ...interface{}) error {
	return NewErrorf(http.StatusNotFound, format, args...)
}

// Conflictf sends a 409 Conflict error with formatted message
func Conflictf(format string, args ...interface{}) error {
	return NewErrorf(http.StatusConflict, format, args...)
}

// InternalErrorf sends a 500 Internal Server Error with formatted message
func InternalErrorf(format string, args ...interface{}) error {
	return NewErrorf(http.StatusInternalServerError, format, args...)
}

// NotImplementedf sends a 501 Not Implemented error with formatted message
func NotImplementedf(format string, args ...interface{}) error {
	return NewErrorf(http.StatusNotImplemented, format, args...)
}

// Unavailablef sends a 503 Service Unavailable error with formatted message
func Unavailablef(format string, args ...interface{}) error {
	return NewErrorf(http.StatusServiceUnavailable, format, args...)
}
