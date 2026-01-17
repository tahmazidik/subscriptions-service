package handler

import (
	"time"

	"github.com/tahmazidik/subscriptions-service/internal/subscription/model"
)

type subscriptionResponse struct {
	ID          string    `json:"id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      string    `json:"user_id"`
	StartDate   string    `json:"start_date"`
	EndDate     *string   `json:"end_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func toResponse(s model.Subscription) subscriptionResponse {
	resp := subscriptionResponse{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		UserID:      s.UserID,
		StartDate:   s.StartDate.Format("01-2006"),
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}
	if s.EndDate != nil {
		end := s.EndDate.Format("01-2006")
		resp.EndDate = &end
	}
	return resp
}
