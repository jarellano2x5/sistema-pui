package usuario

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jarellano2x5/finance-go/internal/utils"
)

var validate = validator.New()

type UsuarioHandler interface {
	PostLogin(c *fiber.Ctx) error
	PostRegister(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
}

type UsuarioHandlerImpl struct {
	srv UsuarioService
}

func NewUsuarioHandler(srv UsuarioService) UsuarioHandler {
	return &UsuarioHandlerImpl{srv}
}

func (h *UsuarioHandlerImpl) PostLogin(c *fiber.Ctx) error {
	req := new(LoginDto)
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if err := validate.Struct(req); err != nil {
		ers := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			ers[err.Field()] = fmt.Sprintf("failed on '%s' rule", err.Tag())
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": ers,
		})
	}
	dto, res := h.srv.Login(*req)
	if res != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to login",
		})
	}
	tk, err := utils.GenerateJWT(dto.Id, dto.User)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"warning": "could not generate token",
		})
	}
	reg := TokenDto{tk}
	return c.Status(fiber.StatusOK).JSON(reg)
}

func (h *UsuarioHandlerImpl) PostRegister(c *fiber.Ctx) error {
	req := new(UsuarioDto)
	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if err := validate.Struct(req); err != nil {
		ers := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			ers[err.Field()] = fmt.Sprintf("failed on '%s' rule", err.Tag())
		}
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": ers,
		})
	}
	dto, res := h.srv.Register(*req)
	if res != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": res.Error(),
		})
	}
	tk, err := utils.GenerateJWT(dto.Id, dto.User)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"warning": "could not generate token",
		})
	}
	reg := TokenDto{tk}
	return c.Status(fiber.StatusCreated).JSON(reg)
}

func (h *UsuarioHandlerImpl) Delete(c *fiber.Ctx) error {
	id, err := utils.GetUserId(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	res := h.srv.Cancel(id)
	if res != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to cancel subscription",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"action": "subscription cancelled!",
	})
}
