package main

import (
	"backend-service/config"
	"backend-service/internal/application"
	"backend-service/internal/infrastructure/api/graph"
	"backend-service/internal/infrastructure/database"
	"backend-service/internal/transport/rest"
	resthandler "backend-service/internal/transport/rest/handler"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	gqlhandler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
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
	db := database.ConnectPg(ctx)
	defer func() {
		if err := db.Close(); err != nil {
			slog.Error("error closing database", "error", err)
		}
	}()

	// Application layer — single composition root
	app := application.New(db)

	// Router
	r := chi.NewRouter()

	// --- Core middlewares ---
	r.Use(middleware.RequestID)  // Adds X-Request-Id to every request
	r.Use(middleware.RealIP)     // Reads X-Forwarded-For / X-Real-IP
	r.Use(middleware.Recoverer)  // Recovers from panics and returns 500
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   cfg.OriginsAllowed,
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
	}).Handler)

	// --- Observability endpoints ---
	r.Get("/health", resthandler.Health)
	r.Get("/readiness", resthandler.Readiness(db))

	// --- REST transport (first-class) ---
	r.Mount("/api/v1", rest.NewRouter(app))

	// --- GraphQL transport (optional adapter) ---
	gqlSrv := gqlhandler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{App: app},
	}))
	r.Handle("/", playground.Handler("GraphQL playground", "/query"))
	r.Handle("/query", gqlSrv)

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
	signal.Notify(quit, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("graceful shutdown failed", "error", err)
	}

	slog.Info("server stopped")
}
