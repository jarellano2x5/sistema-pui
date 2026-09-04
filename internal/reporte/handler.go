package reporte

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	responses "github.com/jarellano2x5/finance-go/internal"
	"github.com/jarellano2x5/finance-go/internal/utils"
)

type ReporteHandler interface {
	GetReportes(c fiber.Ctx) error
	GetGenerados(c fiber.Ctx) error
	PostReporte(c fiber.Ctx) error
	GetConectar(c fiber.Ctx) error
}

type ReporteHandlerImpl struct {
	srv ReporteService
}

func NewReporteHandler(srv ReporteService) ReporteHandler {
	return &ReporteHandlerImpl{srv}
}

func (h *ReporteHandlerImpl) GetReportes(c fiber.Ctx) error {
	id, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ls, err := h.srv.GetAll(id)
	return c.Status(fiber.StatusOK).JSON(ls)
}

func (h *ReporteHandlerImpl) GetGenerados(c fiber.Ctx) error {
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
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"param": err,
		})
	}
	return c.Status(fiber.StatusOK).JSON(ls)
}

func (h *ReporteHandlerImpl) GetConectar(c fiber.Ctx) error {
	id, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"param": err.Error(),
		})
	}
	err = h.srv.SendReport(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.Response{
			Message: err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(responses.Response{
		Message: "Terminado",
	})
}

func (h *ReporteHandlerImpl) PostReporte(c fiber.Ctx) error {
	id, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"param": err.Error(),
		})
	}
	var curps []string
	if err = c.Bind().Body(&curps); err != nil {
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
	return c.Status(fiber.StatusCreated).JSON(responses.Response{
		Message: "Reporte creado",
	})
}
