package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

func ErrorPanic(err error) {
	if err != nil {
		panic(err)
	}
}

func GetUserId(c fiber.Ctx) (uuid.UUID, error) {
	raw := c.Locals("user_id")
	str, ok := raw.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("user_id not found or invalid type")
	}
	id, err := uuid.Parse(str)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid uuid: %w", err)
	}
	return id, nil
}
