package main

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v5/pgxpool"
	log "github.com/sirupsen/logrus"
	"petsitting-api/internal/config"
	"petsitting-api/internal/db"
	"petsitting-api/internal/server/web"
	"petsitting-api/lib/healthcheck"
	"petsitting-api/lib/users"
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
	users, err := users.NewUsers(dbPool, ctx)
	if err != nil {
		log.Fatalf("Failed to construct health check: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
	r.Use(render.SetContentType(render.ContentTypeJSON))
	r.Get("/v1/health", check.GetHealthCheckHandler())
	r.Get("/v1/users", users.GetUsersHandler())
	r.Post("/v1/users/createUser", users.PostUserHandler())
	return r
}
