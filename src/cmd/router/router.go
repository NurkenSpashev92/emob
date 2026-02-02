package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/emob/internal/handler"
	"github.com/nurkenspashev92/emob/internal/initializers"
	"github.com/nurkenspashev92/emob/internal/middleware"
)

func RegisterRoutes(db *pgxpool.Pool) *fiber.App {
	app := fiber.New()

	app.Use(middleware.CorsHandler)
	app.Use(initializers.NewLogger())
	app.Use(initializers.NewSwagger())

	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/healthcheck", handler.HealthCheck(db))
		apiV1.Get("/subscription", handler.GetSubscriptions(db))
	}

	return app
}
