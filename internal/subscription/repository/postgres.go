package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/tahmazidik/subscriptions-service/internal/subscription/model"
)

type Repo struct {
	pool *pgxpool.Pool
}

func NewRepo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

func (r *Repo) Create(ctx context.Context, s model.Subscription) (model.Subscription, error) {
	const q = `
INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at;
`
	var out model.Subscription
	err := r.pool.QueryRow(ctx, q,
		s.ServiceName,
		s.Price,
		s.UserID,
		// важный момент: мы хотим именно "первое число месяца"
		time.Date(s.StartDate.Year(), s.StartDate.Month(), 1, 0, 0, 0, 0, time.UTC),
		normalizeMonthPtr(s.EndDate),
	).Scan(
		&out.ID,
		&out.ServiceName,
		&out.Price,
		&out.UserID,
		&out.StartDate,
		&out.EndDate,
		&out.CreatedAt,
		&out.UpdatedAt,
	)
	return out, err
}

func normalizeMonthPtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	tt := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	return &tt
}
