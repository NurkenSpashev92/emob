package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/emob/internal/handler"
	"github.com/nurkenspashev92/emob/internal/initializers"
	"github.com/nurkenspashev92/emob/internal/middleware"
)

func RegisterRoutes(db *pgxpool.Pool) *fiber.App {
	app := fiber.New(initializers.NewFiberConfig())

	app.Use(middleware.CorsHandler)
	app.Use(initializers.NewLogger())
	app.Use(initializers.NewSwagger())

	apiV1 := app.Group("/api/v1")
	{
		apiV1.Get("/healthcheck", handler.HealthCheck(db))

		apiV1.Get("/subscriptions", handler.GetSubscriptions(db))
		apiV1.Post("/subscriptions", handler.CreateSubscription(db))
		apiV1.Get("/subscriptions/total", handler.GetSubscriptionsTotal(db))
		apiV1.Get("/subscriptions/:id", handler.GetSubscription(db))
		apiV1.Put("/subscriptions/:id", handler.UpdateSubscription(db))
		apiV1.Delete("/subscriptions/:id", handler.DeleteSubscription(db))
	}

	return app
}
