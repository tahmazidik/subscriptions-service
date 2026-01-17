package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/tahmazidik/subscriptions-service/internal/subscription/model"
	"github.com/tahmazidik/subscriptions-service/internal/subscription/repository"
)

type Handler struct {
	repo *repository.Repo
}

func NewHandler(repo *repository.Repo) *Handler {
	return &Handler{repo: repo}
}

type createRequest struct {
	ServiceName string  `json:"service_name"`
	Price       int     `json:"price"`
	UserID      string  `json:"user_id"`
	StartDate   string  `json:"start_date"` // "07-2025"
	EndDate     *string `json:"end_date"`   // "10-2025" or null
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	// простая валидация
	req.ServiceName = strings.TrimSpace(req.ServiceName)
	req.UserID = strings.TrimSpace(req.UserID)

	if req.ServiceName == "" {
		http.Error(w, "service_name is required", http.StatusBadRequest)
		return
	}
	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if req.Price < 0 {
		http.Error(w, "price must be >= 0", http.StatusBadRequest)
		return
	}

	start, err := parseMonthYear(req.StartDate)
	if err != nil {
		http.Error(w, "start_date must be MM-YYYY", http.StatusBadRequest)
		return
	}

	var end *time.Time
	if req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		ed, err := parseMonthYear(*req.EndDate)
		if err != nil {
			http.Error(w, "end_date must be MM-YYYY or null", http.StatusBadRequest)
			return
		}
		end = &ed
	}

	created, err := h.repo.Create(r.Context(), model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   start,
		EndDate:     end,
	})
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(created)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	s, ok, err := h.repo.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(s)
}

func parseMonthYear(s string) (time.Time, error) {
	// формат "MM-YYYY"
	t, err := time.Parse("01-2006", strings.TrimSpace(s))
	if err != nil {
		return time.Time{}, err
	}
	// приводим к первому числу месяца (и UTC)
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC), nil
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	serviceName := strings.TrimSpace(r.URL.Query().Get("service_name"))

	items, err := h.repo.List(r.Context(), userID, serviceName)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json: "+err.Error(), http.StatusBadRequest)
		return
	}

	req.ServiceName = strings.TrimSpace(req.ServiceName)
	req.UserID = strings.TrimSpace(req.UserID)

	if req.ServiceName == "" {
		http.Error(w, "service_name is required", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}

	if req.Price < 0 {
		http.Error(w, "price must be >= 0", http.StatusBadRequest)
		return
	}

	start, err := parseMonthYear(req.StartDate)
	if err != nil {
		http.Error(w, "start_date must be MM-YYYY", http.StatusBadRequest)
		return
	}

	var end *time.Time
	if req.EndDate != nil && strings.TrimSpace(*req.EndDate) != "" {
		ed, err := parseMonthYear(*req.EndDate)
		if err != nil {
			http.Error(w, "end_date must be MM-YYYY or null", http.StatusBadRequest)
			return
		}
		end = &ed
	}

	updated, ok, err := h.repo.Update(r.Context(), id, model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      req.UserID,
		StartDate:   start,
		EndDate:     end,
	})

	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updated)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimSpace(chi.URLParam(r, "id"))
	if id == "" {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}

	ok, err := h.repo.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "subscription not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) Total(w http.ResponseWriter, r *http.Request) {
	userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
	serviceName := strings.TrimSpace(r.URL.Query().Get("service_name"))
	startStr := strings.TrimSpace(r.URL.Query().Get("start_date"))
	endStr := strings.TrimSpace(r.URL.Query().Get("end_date"))

	if userID == "" {
		http.Error(w, "user_id is required", http.StatusBadRequest)
		return
	}
	if startStr == "" || endStr == "" {
		http.Error(w, "start_date and end_date are required", http.StatusBadRequest)
		return
	}

	periodStart, err := parseMonthYear(startStr)
	if err != nil {
		http.Error(w, "start_date must be MM-YYYY", http.StatusBadRequest)
		return
	}
	periodEnd, err := parseMonthYear(endStr)
	if err != nil {
		http.Error(w, "end_date must be MM-YYYY", http.StatusBadRequest)
		return
	}
	if monthIndex(periodStart) > monthIndex(periodEnd) {
		http.Error(w, "start_date must be before or equal to end_date", http.StatusBadRequest)
		return
	}

	subs, err := h.repo.ListForPeriod(r.Context(), userID, serviceName, periodStart, periodEnd)
	if err != nil {
		http.Error(w, "db error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	total := 0
	for _, sub := range subs {
		subEnd := periodEnd
		if sub.EndDate != nil {
			subEnd = *sub.EndDate
		}

		overlapStart := maxMonth(periodStart, sub.StartDate)
		overlapEnd := minMonth(periodEnd, subEnd)

		if monthIndex(overlapStart) > monthIndex(overlapEnd) {
			continue
		}

		months := monthIndex(overlapEnd) - monthIndex(overlapStart) + 1
		total += months * sub.Price
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]int{"total": total})
}

func monthIndex(t time.Time) int {
	return t.Year()*12 + int(t.Month())
}

func maxMonth(a, b time.Time) time.Time {
	if monthIndex(a) >= monthIndex(b) {
		return a
	}
	return b
}

func minMonth(a, b time.Time) time.Time {
	if monthIndex(a) <= monthIndex(b) {
		return a
	}
	return b
}
