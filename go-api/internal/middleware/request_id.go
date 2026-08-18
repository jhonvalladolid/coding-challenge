package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gofiber/fiber/v2"
	"go-api/internal/response"
)

func RequestID() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Get(response.RequestIDHeader)
		if id == "" {
			id = newRequestID()
		}

		c.Locals(response.RequestIDLocalsKey, id)
		c.Set(response.RequestIDHeader, id)
		return c.Next()
	}
}

func newRequestID() string {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "generated-request-id"
	}

	raw[6] = (raw[6] & 0x0f) | 0x40
	raw[8] = (raw[8] & 0x3f) | 0x80

	encoded := hex.EncodeToString(raw[:])
	return encoded[0:8] + "-" + encoded[8:12] + "-" + encoded[12:16] + "-" + encoded[16:20] + "-" + encoded[20:32]
}
