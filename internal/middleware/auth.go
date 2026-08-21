package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jarellano2x5/finance-go/internal/utils"
)

func Protected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		head := c.Get("Authorization")
		if head == "" || !strings.HasPrefix(head, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid token",
			})
		}
		tst := strings.TrimPrefix(head, "Bearer ")
		claims, err := utils.ParseJWT(tst)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing or invalid token",
			})
		}
		c.Locals("user_id", claims["user_id"])
		fmt.Printf("######## %s", claims["user_id"])
		return c.Next()
	}
}
