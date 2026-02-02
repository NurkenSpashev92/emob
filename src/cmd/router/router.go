package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/emob/internal/handler"
)

func RegisterRoutes(db *pgxpool.Pool) *fiber.App {
	app := fiber.New()

	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/healthcheck", handler.HealthCheck(db))
	}

	return app
}
