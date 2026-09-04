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
		option.Summary("Acceso"),
		option.Description("Retorna un token de acceso"),
		option.Request(new(usuario.LoginDto)),
		option.Response(200, new(usuario.TokenDto)),
	)
	r.Post("/notificar-coincidencia", middleware.Protected(), solHand.PostNotificar).With(
		option.Summary("Notificar coincidencia"),
		option.Description("Permite notificar el resulta coincidencias"),
		option.Request(new(solicitud.SolicitudDto)),
		option.Response(200, new(responses.Response)),
		option.Security("bearerAuth"),
	)
	r.Post("/busqueda-finalizada", middleware.Protected(), solHand.PostCerrar).With(
		option.Summary("Búsqueda finalizada"),
		option.Description("Permite notificar el fin de la busqueda"),
		option.Request(new(solicitud.CierreDto)),
		option.Response(200, new(responses.Response)),
		option.Security("bearerAuth"),
	)
	v1 := r.Group("/users")
	{
		v1.Post("/", usuHand.PostRegister).With(
			option.Summary("Registro en esta plataforma"),
			option.Description("Crea usuario y regresa un token"),
			option.Request(new(usuario.UsuarioDto)),
			option.Response(200, new(usuario.TokenDto)),
		)
		v1.Delete("/", middleware.Protected(), usuHand.Delete).With(
			option.Summary("Baja de esta plataforma"),
			option.Description("Permite eliminar al usuario"),
			option.Response(200, new(responses.Response)),
			option.Security("bearerAuth"),
		)
	}
	v1.With(
		option.GroupTags("Registro de usuarios"),
	)
	rp := r.Group("/reportes", middleware.Protected())
	{
		rp.Get("/", repHand.GetReportes).With(
			option.Summary("Lista de reportes"),
			option.Description("Regresa lista de reportes dados de alta"),
			option.Response(200, new([]reporte.ReporteDto)),
			option.Security("bearerAuth"),
		)
		rp.Get("/:cant", repHand.GetGenerados).With(
			option.Summary("Crear registros"),
			option.Description("Genera la cantidad de CURPs aleatorios"),
			option.Response(200, new([]string)),
			option.Security("bearerAuth"),
		)
		rp.Patch("/", repHand.GetConectar).With(
			option.Summary("Envía registros"),
			option.Description("Envía los registro a la entidad diversa"),
			option.Response(200, new([]string)),
			option.Security("bearerAuth"),
		)
		rp.Post("/", repHand.PostReporte).With(
			option.Summary("Registra reportes"),
			option.Description("Registra los CURPs en reportes"),
			option.Request(new([]string)),
			option.Response(200, new(responses.Response)),
			option.Security("bearerAuth"),
		)
	}
	rp.With(
		option.GroupTags("Registro de solicitudes"),
	)
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
