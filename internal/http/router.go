package http

import (
	"context"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	subhandler "github.com/tahmazidik/subscriptions-service/internal/subscription/handler"
	subrepo "github.com/tahmazidik/subscriptions-service/internal/subscription/repository"
)

func NewRouter(pool *pgxpool.Pool) http.Handler {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	r.Get("/db/health", func(w http.ResponseWriter, r *http.Request) {
		if err := pingDB(r.Context(), pool); err != nil {
			http.Error(w, "db not ok: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("db ok"))
	})

	repo := subrepo.NewRepo(pool)
	handler := subhandler.NewHandler(repo)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/subscriptions", handler.Create)
		r.Get("/subscriptions/{id}", handler.GetByID)
	})

	return r
}

func pingDB(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return pool.Ping(ctx)
}
