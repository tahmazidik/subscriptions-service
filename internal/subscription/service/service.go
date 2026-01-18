package service

import (
	"context"
	"time"

	"github.com/tahmazidik/subscriptions-service/internal/subscription/model"
)

type Repository interface {
	ListForPeriod(ctx context.Context, userID, serviceName string, periodStart, periodEnd time.Time) ([]model.Subscription, error)
}

type Service struct {
	repo Repository
}

func New(repo Repository) *Service {
	return &Service{repo: repo}
}
