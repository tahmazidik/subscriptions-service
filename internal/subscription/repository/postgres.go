package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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

func (r *Repo) GetByID(ctx context.Context, id string) (model.Subscription, bool, error) {
	const q = `
SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
FROM subscriptions
WHERE id = $1;
`
	var out model.Subscription
	err := r.pool.QueryRow(ctx, q, id).Scan(
		&out.ID,
		&out.ServiceName,
		&out.Price,
		&out.UserID,
		&out.StartDate,
		&out.EndDate,
		&out.CreatedAt,
		&out.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscription{}, false, nil
		}
		return model.Subscription{}, false, err
	}

	return out, true, nil
}

func (r *Repo) List(ctx context.Context, userID, serviceName string) ([]model.Subscription, error) {
	const q = `
SELECT id, service_name, price, user_id, start_date, end_date, created_at, updated_at
FROM subscriptions
WHERE ($1 = '' OR user_id::text = $1)
  AND ($2 = '' OR service_name = $2)
ORDER BY created_at DESC;
`
	rows, err := r.pool.Query(ctx, q, userID, serviceName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]model.Subscription, 0)
	for rows.Next() {
		var s model.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.ServiceName,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *Repo) Update(ctx context.Context, id string, s model.Subscription) (model.Subscription, bool, error) {
	const q = `
UPDATE subscriptions
SET service_name = $2,
    price = $3,
    user_id = $4,
    start_date = $5,
    end_date = $6,
    updated_at = now()
WHERE id = $1
RETURNING id, service_name, price, user_id, start_date, end_date, created_at, updated_at;
`
	var out model.Subscription
	err := r.pool.QueryRow(ctx, q,
		id,
		s.ServiceName,
		s.Price,
		s.UserID,
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

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Subscription{}, false, nil
		}
		return model.Subscription{}, false, err
	}

	return out, true, nil
}

func (r *Repo) Delete(ctx context.Context, id string) (bool, error) {
	const q = `
DELETE FROM subscriptions
WHERE id = $1;
`
	cmdTag, err := r.pool.Exec(ctx, q, id)
	if err != nil {
		return false, err
	}
	return cmdTag.RowsAffected() > 0, nil
}

func normalizeMonthPtr(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	tt := time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
	return &tt
}
