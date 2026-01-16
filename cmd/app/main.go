package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"

	"github.com/tahmazidik/subscriptions-service/internal/config"
	"github.com/tahmazidik/subscriptions-service/internal/db"
)

func main() {
	_ = godotenv.Load()

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := db.NewPool(ctx, cfg.DSN())
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	defer pool.Close()

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

	port := cfg.AppPort
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func pingDB(ctx context.Context, pool *pgxpool.Pool) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	return pool.Ping(ctx)
}
