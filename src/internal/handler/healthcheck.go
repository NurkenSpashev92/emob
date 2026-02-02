package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

// HealthCheck godoc
// @Summary      Health Check
// @Description  Checks if the application and database are running
// @Tags         Health
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string{"status":"ok","message":"success"}
// @Failure      503  {object}  map[string]string{"status":"fail","message":"database not reachable"}
// @Router       /health [get]
func HealthCheck(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if err := db.Ping(c.Context()); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "fail",
				"message": "database not reachable",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "ok",
			"message": "success",
		})
	}
}
