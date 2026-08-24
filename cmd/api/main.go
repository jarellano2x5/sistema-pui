package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/jarellano2x5/finance-go/config"
	"github.com/jarellano2x5/finance-go/database"
	responses "github.com/jarellano2x5/finance-go/internal"
	"github.com/jarellano2x5/finance-go/internal/middleware"
	"github.com/jarellano2x5/finance-go/internal/reporte"
	"github.com/jarellano2x5/finance-go/internal/solicitud"
	"github.com/jarellano2x5/finance-go/internal/usuario"
	"github.com/oaswrap/spec/adapter/fiberv3openapi"
	"github.com/oaswrap/spec/option"
)

func main() {
	// configuration loading
	cfg := config.LoadConfig()
	// database connection
	database.Connect(cfg)

	// define fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Test PUI API v1.0",
		ErrorHandler: customErrorHandler,
	})
	// OpenAPI router
	r := fiberv3openapi.NewRouter(app,
		option.WithTitle("Test PUI API v1.0"),
		option.WithVersion("1.0.0"),
		option.WithSecurity("bearerAuth", option.SecurityHTTPBearer("Bearer")),
	)
	// middleware initialization
	middleware.SetupMiddleware(app)
	// checking health
	r.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "test-pui-api",
		})
	}).With(
		option.Summary("Salud del aplicativo"),
	)
	// implements
	usuSer := usuario.NewUsuarioService(database.DB)
	usuHand := usuario.NewUsuarioHandler(usuSer)
	repSer := reporte.NewReporteService(database.DB)
	repHand := reporte.NewReporteHandler(repSer)
	solSer := solicitud.NewSolicitudService(database.DB)
	solHand := solicitud.NewSolicitudHandler(solSer)
	// define routes
	r.Post("/login", usuHand.PostLogin).With(
		option.Request(new(usuario.LoginDto)),
		option.Response(200, new(usuario.TokenDto)),
	)
	r.Post("/notificar-coincidencia", middleware.Protected(), solHand.PostNotificar).With(
		option.Request(new(solicitud.SolicitudDto)),
		option.Response(200, new(responses.Response)),
	)
	r.Post("/busqueda-finalizada", middleware.Protected(), solHand.PostCerrar).With(
		option.Request(new(solicitud.CierreDto)),
		option.Response(200, new(responses.Response)),
	)
	v1 := r.Group("/users")
	{
		v1.Post("/", usuHand.PostRegister)
		v1.Delete("/", middleware.Protected(), usuHand.Delete)
	}
	rp := r.Group("/reportes", middleware.Protected())
	{
		rp.Get("/", repHand.GetReportes)
		rp.Get("/:cant", repHand.GetGenerados)
		rp.Post("/", repHand.PostReporte)
	}
	// server init
	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("Server starting on %s", addr)
	log.Fatal(app.Listen(addr))
}

func customErrorHandler(c fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}
	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}
