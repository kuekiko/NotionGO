package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 表示错误代码
type ErrorCode string

const (
	// API 错误代码
	ErrValidation          ErrorCode = "validation_error"
	ErrInvalidRequest      ErrorCode = "invalid_request"
	ErrInvalidJSON         ErrorCode = "invalid_json"
	ErrMissingVersion      ErrorCode = "missing_version"
	ErrUnauthorized        ErrorCode = "unauthorized"
	ErrRateLimited         ErrorCode = "rate_limited"
	ErrInternalServerError ErrorCode = "internal_server_error"
	ErrConflict            ErrorCode = "conflict"
	ErrServiceUnavailable  ErrorCode = "service_unavailable"

	// SDK 错误代码
	ErrInvalidInput      ErrorCode = "invalid_input"
	ErrSizeLimitExceeded ErrorCode = "size_limit_exceeded"
	ErrRequestTimeout    ErrorCode = "request_timeout"
	ErrContextCanceled   ErrorCode = "context_canceled"
)

// Error 表示 Notion API 错误
type Error struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Status     int       `json:"status"`
	RetryAfter int       `json:"retry_after,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("notion: %s (status: %d, code: %s)", e.Message, e.Status, e.Code)
}

// NewError 创建一个新的错误
func NewError(code ErrorCode, message string, status int) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// IsNotFound 检查是否是 404 错误
func IsNotFound(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Status == http.StatusNotFound
	}
	return false
}

// IsRateLimited 检查是否是速率限制错误
func IsRateLimited(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrRateLimited
	}
	return false
}

// IsSizeLimitExceeded 检查是否超出大小限制
func IsSizeLimitExceeded(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrSizeLimitExceeded
	}
	return false
}

// IsValidationError 检查是否是验证错误
func IsValidationError(err error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == ErrValidation
	}
	return false
}

// SizeLimits 定义各种大小限制
var SizeLimits = struct {
	MaxRichTextContent    int
	MaxRichTextLinkURL    int
	MaxEquationExpression int
	MaxArrayElements      int
	MaxURL                int
	MaxEmail              int
	MaxPhoneNumber        int
	MaxMultiSelect        int
	MaxRelation           int
	MaxPeople             int
	MaxPayloadBlocks      int
	MaxPayloadSize        int64 // 字节
}{
	MaxRichTextContent:    2000,
	MaxRichTextLinkURL:    2000,
	MaxEquationExpression: 1000,
	MaxArrayElements:      100,
	MaxURL:                2000,
	MaxEmail:              200,
	MaxPhoneNumber:        200,
	MaxMultiSelect:        100,
	MaxRelation:           100,
	MaxPeople:             100,
	MaxPayloadBlocks:      1000,
	MaxPayloadSize:        512 * 1024, // 500KB
}
