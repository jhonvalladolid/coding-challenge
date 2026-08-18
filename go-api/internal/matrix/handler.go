package matrix

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go-api/internal/config"
	apperr "go-api/internal/errors"
	"go-api/internal/middleware"
	"go-api/internal/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func NewApp(cfg config.Config) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler:          middleware.ErrorHandler,
		BodyLimit:             1024 * 1024,
		DisableStartupMessage: cfg.AppEnv == "test",
	})

	app.Use(recover.New())
	app.Use(middleware.RequestID())

	handler := NewHandler(NewService(cfg.MaxMatrixDim))
	app.Get("/health", handler.Health)
	app.Post("/api/v1/matrices/qr", handler.Factorize)

	return app
}

func (h *Handler) Health(c *fiber.Ctx) error {
	return response.Success(c, HealthData{
		Service: "go-api",
		Status:  "ok",
	})
}

func (h *Handler) Factorize(c *fiber.Ctx) error {
	var req QRRequest
	if err := c.BodyParser(&req); err != nil {
		return apperr.InvalidRequestBody()
	}

	result, err := h.service.Factorize(req.Matrix)
	if err != nil {
		return err
	}

	return response.Success(c, result)
}
