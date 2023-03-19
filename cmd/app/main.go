package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
	"golang-starter/internal/config"
	"golang-starter/internal/db"
	"golang-starter/internal/server/web"
	"golang-starter/lib/healthcheck"
)

func main() {
	ctx := context.Background()
	configuration := config.Init()

	dbPool := db.GetDbConnectionFactory(configuration.DbConfig.Url)

	server := web.NewServer(configuration.ServerConfig.Address)
	server.AddRoute("/api", setupRoutes(dbPool, ctx))

	log.Infof("Starting server on %s...", server.Address)
	server.Start(ctx)
}

func setupRoutes(dbPool *pgxpool.Pool, ctx context.Context) chi.Router {
	check, err := healthcheck.NewHealthCheck(dbPool, ctx)
	if err != nil {
		log.Fatalf("Failed to construct health check: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/v1/health", check.GetHealthCheckHandler())

	return r
}
