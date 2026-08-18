package errors

import stderrors "errors"

type AppError struct {
	Code       string
	Message    string
	StatusCode int
}

func (e *AppError) Error() string {
	return e.Message
}

func AsAppError(err error) (*AppError, bool) {
	var appErr *AppError
	if stderrors.As(err, &appErr) {
		return appErr, true
	}
	return nil, false
}

func InvalidRequestBody() *AppError {
	return &AppError{
		Code:       "INVALID_REQUEST_BODY",
		Message:    "The request body is invalid or is not valid JSON.",
		StatusCode: 400,
	}
}

func MatrixRequired() *AppError {
	return &AppError{
		Code:       "MATRIX_REQUIRED",
		Message:    "The property matrix is required.",
		StatusCode: 400,
	}
}

func EmptyMatrix() *AppError {
	return &AppError{
		Code:       "EMPTY_MATRIX",
		Message:    "The matrix must contain at least one row.",
		StatusCode: 400,
	}
}

func EmptyMatrixRow() *AppError {
	return &AppError{
		Code:       "EMPTY_MATRIX_ROW",
		Message:    "Each matrix row must contain at least one column.",
		StatusCode: 400,
	}
}

func IrregularMatrix() *AppError {
	return &AppError{
		Code:       "IRREGULAR_MATRIX",
		Message:    "All matrix rows must contain the same number of elements.",
		StatusCode: 400,
	}
}

func InvalidMatrixValue() *AppError {
	return &AppError{
		Code:       "INVALID_MATRIX_VALUE",
		Message:    "All matrix values must be finite numbers.",
		StatusCode: 400,
	}
}

func UnsupportedMatrixDimensions(message string) *AppError {
	return &AppError{
		Code:       "UNSUPPORTED_MATRIX_DIMENSIONS",
		Message:    message,
		StatusCode: 422,
	}
}

func QRFactorizationError() *AppError {
	return &AppError{
		Code:       "QR_FACTORIZATION_ERROR",
		Message:    "The QR factorization could not be completed.",
		StatusCode: 422,
	}
}

func InternalServerError() *AppError {
	return &AppError{
		Code:       "INTERNAL_SERVER_ERROR",
		Message:    "An unexpected error occurred.",
		StatusCode: 500,
	}
}

func NotFound() *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    "The requested resource was not found.",
		StatusCode: 404,
	}
}
