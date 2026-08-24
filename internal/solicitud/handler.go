package solicitud

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate = validator.New()

type SolicitudHandler interface {
	PostNotificar(c fiber.Ctx) error
	PostCerrar(c fiber.Ctx) error
}

type SolicitudHandlerImpl struct {
	srv SolicitudService
}

func NewSolicitudHandler(srv SolicitudService) SolicitudHandler {
	return &SolicitudHandlerImpl{srv}
}

func (h *SolicitudHandlerImpl) PostNotificar(c fiber.Ctx) error {
	sol := new(SolicitudDto)
	if err := c.Bind().Body(&sol); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := validate.Struct(sol); err != nil {
		ers := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			ers[err.Field()] = fmt.Sprintf("failed on '%s' rule", err.Tag())
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": ers,
		})
	}
	if err := h.srv.Notificar(*sol); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Coincidencia recibida correctamente",
	})
}

func (h *SolicitudHandlerImpl) PostCerrar(c fiber.Ctx) error {
	var cie = new(CierreDto)
	if err := c.Bind().Body(&cie); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	if err := validate.Struct(cie); err != nil {
		ers := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			ers[err.Field()] = fmt.Sprintf("failed on '%s' rule", err.Tag())
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": ers,
		})
	}
	if err := h.srv.Cierre(*cie); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Coincidencia recibida correctamente",
	})
}
