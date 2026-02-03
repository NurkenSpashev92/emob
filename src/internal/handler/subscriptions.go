package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/emob/internal/models"
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
// @Router       /api/v1/subscriptions [get]
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
			fmt.Println("Error:", err)
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.JSON(subscriptions)
	}
}

// CreateSubscription godoc
// @Summary      Create subscription
// @Description  Create a new subscription
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        body  body      models.CreateSubscription  true  "Subscription body"
// @Success      201   {object}  models.Subscription
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /api/v1/subscriptions [post]
func CreateSubscription(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var body models.CreateSubscription

		if err := c.BodyParser(&body); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}

		repo := repositories.NewSubscriptionRepository(db)

		subscription, err := repo.CreateSubscriptions(c.Context(), body)
		if err != nil {
			fmt.Println("Error:", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(subscription)
	}
}

// GetSubscription godoc
// @Summary      Get subscription by ID
// @Description  Returns a single subscription
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Subscription ID"
// @Success      200  {object}  models.Subscription
// @Failure      404  {object}  interface{}
// @Router       /api/v1/subscriptions/{id} [get]
func GetSubscription(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		repo := repositories.NewSubscriptionRepository(db)

		sub, err := repo.GetSubscriptionByID(c.Context(), id)
		if err != nil {
			fmt.Println("Error:", err)
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Subscription not found",
			})
		}

		return c.JSON(sub)
	}
}

// UpdateSubscription godoc
// @Summary      Update subscription
// @Description  Updates subscription by ID
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        id    path      string  true  "Subscription ID"
// @Param        body  body      models.CreateSubscription  true  "Subscription body"
// @Success      200   {object}  models.Subscription
// @Failure      400   {object}  interface{}
// @Failure      500   {object}  interface{}
// @Router       /api/v1/subscriptions/{id} [put]
func UpdateSubscription(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		var body models.CreateSubscription
		if err := c.BodyParser(&body); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return c.Status(400).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}

		repo := repositories.NewSubscriptionRepository(db)
		sub, err := repo.UpdateSubscription(c.Context(), id, body)
		if err != nil {
			fmt.Println("Error:", err)
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.JSON(sub)
	}
}

// DeleteSubscription godoc
// @Summary      Delete subscription
// @Description  Deletes subscription by ID
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Subscription ID"
// @Success      204  "No Content"
// @Failure      404  {object}  interface{}
// @Router       /api/v1/subscriptions/{id} [delete]
func DeleteSubscription(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id")
		repo := repositories.NewSubscriptionRepository(db)

		err := repo.DeleteSubscription(c.Context(), id)
		if err != nil {
			fmt.Println("Error:", err)
			return c.Status(404).JSON(fiber.Map{
				"status":  "error",
				"message": "Subscription not found",
			})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

// GetSubscriptionsTotal godoc
// @Summary      Get total subscriptions cost
// @Description  Returns total cost of subscriptions for selected period with optional filters
// @Tags         Subscriptions
// @Accept       json
// @Produce      json
// @Param        date_from    query     string  true   "Start date (YYYY-MM-DD)"
// @Param        date_to      query     string  true   "End date (YYYY-MM-DD)"
// @Param        user_id      query     string  false  "User ID"
// @Param        service_name query     string  false  "Service name"
// @Success      200          {object}  map[string]float64
// @Failure      400          {object}  map[string]string
// @Failure      500          {object}  map[string]string
// @Router       /api/v1/subscriptions/total [get]
func GetSubscriptionsTotal(db *pgxpool.Pool) fiber.Handler {
	return func(c *fiber.Ctx) error {
		dateFrom := c.Query("date_from")
		dateTo := c.Query("date_to")

		if dateFrom == "" || dateTo == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "date_from и date_to обязательны",
			})
		}

		userID := c.Query("user_id")
		serviceName := c.Query("service_name")

		repo := repositories.NewSubscriptionRepository(db)

		total, err := repo.GetTotalSubscriptionsCost(
			c.Context(),
			dateFrom,
			dateTo,
			userID,
			serviceName,
		)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"total_price": total,
		})
	}
}
