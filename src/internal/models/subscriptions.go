package models

// Subscription represents a subscription entity
type Subscription struct {
	ID          string `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	ServiceName string `json:"service_name" example:"Netflix"`
	Price       int    `json:"price" example:"1999"`
	UserID      string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	StartDate   string `json:"start_date" example:"2026-01-01"`
	EndDate     string `json:"end_date,omitempty" example:"2026-12-31"`
	CreatedAt   string `json:"created_at" example:"2026-01-01T12:00:00Z"`
}
