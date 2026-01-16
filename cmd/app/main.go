package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"

	"github.com/tahmazidik/subscriptions-service/internal/config"
	"github.com/tahmazidik/subscriptions-service/internal/db"
	httpapi "github.com/tahmazidik/subscriptions-service/internal/http"
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

	handler := httpapi.NewRouter(pool)

	srv := &http.Server{
		Addr:              ":" + cfg.AppPort,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("listening on :%s", cfg.AppPort)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
