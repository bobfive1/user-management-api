package error

import "time"

type baseError struct {
	ErrorCode    string
	ErrorMessage string
	ErrorDetail  any
}

type InternalServerError struct {
	baseError
}

type FieldValidationError struct {
	baseError
}

type ApiServerError struct {
	baseError
}

type NotFoundError struct {
	baseError
}

type ErrorResponse struct {
	TimeStamp    time.Time `json:"timestamp"`
	ErrorCode    string    `json:"error_code"`
	ErrorMessage string    `json:"error_message"`
	ErrorDetail  any       `json:"error_detail,omitempty"`
}

func (e InternalServerError) Error() string {
	return e.ErrorMessage
}

func NewInternalServerError(detail any) InternalServerError {
	return InternalServerError{
		baseError{
			ErrorCode:    "500",
			ErrorMessage: "Internal server error",
			ErrorDetail:  detail,
		},
	}
}

func (e FieldValidationError) Error() string {
	return e.ErrorMessage
}

func NewFieldValidationError(detail any) FieldValidationError {
	return FieldValidationError{
		baseError{
			ErrorCode:    "400",
			ErrorMessage: "Field validation Error",
			ErrorDetail:  detail,
		},
	}
}

func NewBadRequestError(code, message string, detail any) FieldValidationError {
	return FieldValidationError{
		baseError{
			ErrorCode:    code,
			ErrorMessage: message,
			ErrorDetail:  detail,
		},
	}
}

func (e ApiServerError) Error() string {
	return e.ErrorMessage
}

func NewApiServerError(code, message string, detail any) ApiServerError {
	return ApiServerError{
		baseError{
			ErrorCode:    code,
			ErrorMessage: message,
			ErrorDetail:  detail,
		},
	}
}

func (e NotFoundError) Error() string {
	return e.ErrorMessage
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{
		baseError{
			ErrorCode:    "404",
			ErrorMessage: message,
		},
	}
}
