package fiber

import (
	"context"
	"fmt"
	"time"

	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/common"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/http"
	"github.com/dmitrii-a/hw_go/hw12_13_14_15_calendar/internal/presentation/http/fiber/handlers"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/pprof"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

// NewServer returns a new instance of a server.
func NewServer() http.Server {
	return &server{}
}

type server struct {
	app *fiber.App
}

func (s *server) getServerAddr() string {
	return fmt.Sprintf("%v:%v", common.Config.Server.Host, common.Config.Server.Port)
}

// Start starts the HTTP server.
func (s *server) Start(ctx context.Context) error {
	app := fiber.New()
	s.app = app
	app.Use(
		recover.New(
			recover.Config{
				EnableStackTrace: true,
			},
		),
	)
	if common.Config.Server.Debug {
		app.Use(logger.New(logger.Config{
			Format:     "${time} [${status}] ${latency} ${ip} ${method} ${ua} ${host}${url}\n",
			TimeFormat: time.RFC3339,
			TimeZone:   "Local",
		}))
		app.Use(pprof.New())
	}

	// Routes
	app.Get("/", handlers.HelloWorld)
	api := app.Group("/api/v1")
	api.Get("/health/", handlers.HealthCheck)
	err := app.Listen(s.getServerAddr())
	if common.IsErr(err) {
		return err
	}
	<-ctx.Done()
	return nil
}

// Stop stops the HTTP server.
func (s *server) Stop(ctx context.Context) error {
	common.Logger.Info().Msg("calendar service is stopping...")
	return s.app.ShutdownWithContext(ctx)
}
