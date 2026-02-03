package repositories

import (
	"context"
	"fmt"
	"time"

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

func (repo *SubscriptionRepository) CreateSubscriptions(
	ctx context.Context,
	subscriptionBody models.CreateSubscription,
) (*models.Subscription, error) {

	startDate, err := time.Parse("2006-01-02", subscriptionBody.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}

	var endDate time.Time
	if subscriptionBody.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", subscriptionBody.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end date: %w", err)
		}
	}

	query := `
		INSERT INTO subscriptions (
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at
		) VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at;
	`

	row := repo.db.QueryRow(ctx, query,
		subscriptionBody.ServiceName,
		subscriptionBody.Price,
		subscriptionBody.UserID,
		startDate,
		endDate,
	)

	var s models.Subscription
	if err := row.Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to scan created subscription: %w", err)
	}

	return &s, nil
}

func (repo *SubscriptionRepository) GetSubscriptionByID(
	ctx context.Context,
	id string,
) (*models.Subscription, error) {

	query := `
		SELECT id, service_name, price, user_id, start_date, end_date, created_at
		FROM subscriptions
		WHERE id = $1;
	`

	row := repo.db.QueryRow(ctx, query, id)

	var s models.Subscription
	if err := row.Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &s, nil
}

func (repo *SubscriptionRepository) UpdateSubscription(
	ctx context.Context,
	id string,
	body models.CreateSubscription,
) (*models.Subscription, error) {

	query := `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4, end_date = $5
		WHERE id = $6
		RETURNING id, service_name, price, user_id, start_date, end_date, created_at;
	`

	row := repo.db.QueryRow(ctx, query,
		body.ServiceName,
		body.Price,
		body.UserID,
		body.StartDate,
		body.EndDate,
		id,
	)

	var s models.Subscription
	if err := row.Scan(
		&s.ID,
		&s.ServiceName,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
		&s.CreatedAt,
	); err != nil {
		return nil, fmt.Errorf("failed to scan updated subscription: %w", err)
	}

	return &s, nil
}

func (repo *SubscriptionRepository) DeleteSubscription(
	ctx context.Context,
	id string,
) error {

	query := `DELETE FROM subscriptions WHERE id = $1;`

	_, err := repo.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

func (repo *SubscriptionRepository) GetTotalSubscriptionsCost(
	ctx context.Context,
	dateFrom string,
	dateTo string,
	userID string,
	serviceName string,
) (float64, error) {

	query := `
		SELECT COALESCE(SUM(price), 0)
		FROM subscriptions
		WHERE start_date >= $1
		  AND start_date <= $2
	`

	args := []interface{}{dateFrom, dateTo}
	argID := 3

	if userID != "" {
		query += fmt.Sprintf(" AND user_id = $%d", argID)
		args = append(args, userID)
		argID++
	}

	if serviceName != "" {
		query += fmt.Sprintf(" AND service_name ILIKE $%d", argID)
		args = append(args, serviceName)
	}

	fmt.Println(query)
	var total float64
	err := repo.db.QueryRow(ctx, query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total subscriptions cost: %w", err)
	}

	return total, nil
}
