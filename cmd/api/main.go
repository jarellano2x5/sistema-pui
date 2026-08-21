package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/jarellano2x5/finance-go/config"
	"github.com/jarellano2x5/finance-go/database"
	"github.com/jarellano2x5/finance-go/internal/middleware"
	"github.com/jarellano2x5/finance-go/internal/reporte"
	"github.com/jarellano2x5/finance-go/internal/usuario"
)

func main() {
	// configuration loading
	cfg := config.LoadConfig()
	// database connection
	database.Connect(cfg)

	// define fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Finance API v1.0",
		ErrorHandler: customErrorHandler,
	})
	// middleware initialization
	middleware.SetupMiddleware(app)
	// checking health
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "finance-api",
		})
	})
	// implements
	usuSer := usuario.NewUsuarioService(database.DB)
	usuHand := usuario.NewUsuarioHandler(usuSer)
	repSer := reporte.NewReporteService(database.DB)
	repHand := reporte.NewReporteHandler(repSer)
	// define routes
	v1 := app.Group("/api")
	{
		v1.Post("/login", usuHand.PostLogin)
		users := v1.Group("/users")
		users.Post("/register", usuHand.PostRegister)
		users.Delete("/cancel", middleware.Protected(), usuHand.Delete)
		repos := v1.Group("/reportes", middleware.Protected())
		repos.Get("/", repHand.GetReportes)
		repos.Get("/generar/:cant", repHand.GetGenerados)
		repos.Post("/", repHand.PostReporte)

	}
	// server init
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on %s", addr)
	log.Fatal(app.Listen(addr))
}

func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
