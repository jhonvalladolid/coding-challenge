package response

import (
	"bytes"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
)

const RequestIDLocalsKey = "requestID"
const RequestIDHeader = "X-Request-ID"

func Marshal(v any) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}

	return bytes.TrimRight(buf.Bytes(), "\n"), nil
}

type Meta struct {
	RequestID string `json:"requestId"`
}

type SuccessBody struct {
	Success bool `json:"success"`
	Data    any  `json:"data"`
	Meta    Meta `json:"meta"`
}

type ErrorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorBody struct {
	Success bool         `json:"success"`
	Error   ErrorPayload `json:"error"`
	Meta    Meta         `json:"meta"`
}

func Success(c *fiber.Ctx, data any) error {
	return c.Status(fiber.StatusOK).JSON(SuccessBody{
		Success: true,
		Data:    data,
		Meta:    metaFrom(c),
	})
}

func Error(c *fiber.Ctx, status int, code, message string) error {
	return c.Status(status).JSON(ErrorBody{
		Success: false,
		Error: ErrorPayload{
			Code:    code,
			Message: message,
		},
		Meta: metaFrom(c),
	})
}

func metaFrom(c *fiber.Ctx) Meta {
	requestID, _ := c.Locals(RequestIDLocalsKey).(string)
	return Meta{RequestID: requestID}
}
