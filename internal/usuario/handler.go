package usuario

import "github.com/gofiber/fiber/v2"

type UsuarioHandler interface {
	GetUsuario(c *fiber.Ctx) error
	PostUsuario(c *fiber.Ctx) error
	PutUsuario(c *fiber.Ctx) error
}

type UsuarioHandlerImpl struct {
	srv UsuarioService
}

func NewUsuarioHandler(ser UsuarioService) UsuarioHandler {
	return &UsuarioHandlerImpl{srv: ser}
}

func (h *UsuarioHandlerImpl) GetUsuario(c *fiber.Ctx) error {
	lo, err := h.srv.GetAll()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "something goes wrong",
		})
	}
	return c.Status(fiber.StatusOK).JSON(lo)
}

func (h *UsuarioHandlerImpl) PostUsuario(c *fiber.Ctx) error {
	req := new(UsuarioDto)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if errors := req.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"errors": errors,
		})
	}
	res := h.srv.CreateUser(*req)
	if res != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create usuario",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(req)
}

func (h *UsuarioHandlerImpl) PutUsuario(c *fiber.Ctx) error {
	req := new(UsuarioDto)

	if err := c.BodyParser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	res := h.srv.RepairUser(*req)
	if res != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to udpate usuario",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(req)
}
