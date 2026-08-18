package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	apperr "go-api/internal/errors"
	"go-api/internal/response"
)

func ErrorHandler(c *fiber.Ctx, err error) error {
	if appErr, ok := apperr.AsAppError(err); ok {
		return response.Error(c, appErr.StatusCode, appErr.Code, appErr.Message)
	}

	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		switch fiberErr.Code {
		case fiber.StatusNotFound:
			notFound := apperr.NotFound()
			return response.Error(c, notFound.StatusCode, notFound.Code, notFound.Message)
		case fiber.StatusBadRequest, fiber.StatusUnprocessableEntity:
			invalid := apperr.InvalidRequestBody()
			return response.Error(c, invalid.StatusCode, invalid.Code, invalid.Message)
		}
	}

	internal := apperr.InternalServerError()
	return response.Error(c, internal.StatusCode, internal.Code, internal.Message)
}
