package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/emob/internal/repositories"
)

// GetSubscriptions godoc
// @Summary      Get subscriptions
// @Description  Returns list of subscriptions with pagination
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        limit   query     int  false  "Limit"   default(10)
// @Param        offset query     int  false  "Offset"  default(0)
// @Success      200     {array}   models.Subscription
// @Failure      500     {object}  interface{}
// @Router       /subscription [get]
func GetSubscriptions(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		limit := c.QueryInt("limit", 10)
		offset := c.QueryInt("offset", 0)

		repo := repositories.NewSubscriptionRepository(db)
		subscriptions, err := repo.GetAllSubscriptions(
			c.Context(),
			limit,
			offset,
		)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.JSON(subscriptions)
	}
}
