package reporte

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ReporteHandler interface {
	GetReportes(c *fiber.Ctx) error
	GetGenerados(c *fiber.Ctx) error
	PostReporte(c *fiber.Ctx) error
}

type ReporteHandlerImpl struct {
	srv ReporteService
}

func NewReporteHandler(srv ReporteService) ReporteHandler {
	return &ReporteHandlerImpl{srv}
}

func (h *ReporteHandlerImpl) GetReportes(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"param": "id is required",
		})
	}
	ls, err := h.srv.GetAll(id)
	return c.Status(fiber.StatusOK).JSON(ls)
}

func (h *ReporteHandlerImpl) GetGenerados(c *fiber.Ctx) error {
	qn, err := strconv.Atoi(c.Params("cant"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"param": "cantidad is required",
		})
	}
	if qn > 5 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"param": "cantidad maximum 5",
		})
	}
	ls, err := h.srv.Generate(qn)
	return c.Status(fiber.StatusOK).JSON(ls)
}

func (h *ReporteHandlerImpl) PostReporte(c *fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"param": "id is required",
		})
	}
	var curps []string
	if err = c.BodyParser(&curps); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"body": "curp lists",
		})
	}
	err = h.srv.CreateReport(id, curps)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"body": err,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
	})
}
