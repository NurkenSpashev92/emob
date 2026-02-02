package repositories

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/nurkenspashev92/emob/internal/models"
)

type SubscriptionRepository struct {
	db *pgxpool.Pool
}

func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (repo *SubscriptionRepository) GetAllSubscriptions(
	ctx context.Context,
	limit int,
	offset int,
) ([]models.Subscription, error) {

	query := `
		SELECT 
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at
		FROM subscriptions
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2;
	`

	rows, err := repo.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	subscriptions := make([]models.Subscription, 0)

	for rows.Next() {
		var s models.Subscription

		err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
			&s.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}

		subscriptions = append(subscriptions, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return subscriptions, nil
}
