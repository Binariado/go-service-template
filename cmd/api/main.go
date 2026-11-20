package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/binariado/go-service-template/config"
	"github.com/binariado/go-service-template/internal/application"
	"github.com/binariado/go-service-template/internal/infrastructure/database"
	"github.com/binariado/go-service-template/internal/infrastructure/database/pg"
	graph "github.com/binariado/go-service-template/internal/transport/graphql"
	"github.com/binariado/go-service-template/internal/transport/rest"
	resthandler "github.com/binariado/go-service-template/internal/transport/rest/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func init() {
	// .env is loaded in development. In production, variables come from the
	// environment directly and the missing file is silently ignored.
	_ = godotenv.Load()
}

func main() {
	cfg := config.Load()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Infrastructure
	db := database.Connect(ctx, cfg.PostgresURL)
	defer func() {
		if err := db.Postgres.Close(); err != nil {
			slog.Error("error closing database", "error", err)
		}
	}()
	pg.RunMigrations(db.Postgres)

	// Application layer — single composition root
	app := application.New(db)

	// Router
	r := chi.NewRouter()

	// --- Core middlewares ---
	r.Use(middleware.RequestID) // Adds X-Request-Id to every request
	r.Use(middleware.RealIP)    // Reads X-Forwarded-For / X-Real-IP
	r.Use(middleware.Recoverer) // Recovers from panics and returns 500
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   cfg.OriginsAllowed,
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
	}).Handler)

	// --- Observability endpoints ---
	r.Get("/health", resthandler.Health)
	r.Get("/readiness", resthandler.Readiness(db.Postgres))

	// --- REST transport (first-class) ---
	r.Mount("/api/v1", rest.NewRouter(app))

	// --- GraphQL transport (optional adapter — remove this line to drop it) ---
	graph.RegisterRoutes(r, app, cfg.GraphQLPlaygroundEnabled)

	// --- HTTP server with production-ready timeouts ---
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in background
	go func() {
		slog.Info("server started", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	// --- Graceful shutdown ---
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}

	slog.Info("server stopped")
}
