package healthcheck

import (
	"context"
	"github.com/go-chi/render"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type HealthCheck struct {
	conn *pgxpool.Conn
}

type Result struct {
	DbConnected bool `json:"dbConnected"`
}

func (result *Result) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewHealthCheck(pool *pgxpool.Pool, ctx context.Context) (HealthCheck, error) {
	acquire, err := pool.Acquire(ctx)

	if err != nil {
		return HealthCheck{}, err
	}

	return HealthCheck{
		conn: acquire,
	}, nil
}

func (h HealthCheck) GetHealthCheckHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		checks, healthy := h.doChecks()

		if healthy {
			w.WriteHeader(200)
		} else {
			w.WriteHeader(400)
		}
		render.Render(w, r, checks)
	}
}

func (h HealthCheck) doChecks() (*Result, bool) {
	err := h.conn.Ping(context.Background())
	if err != nil {
		return &Result{
			DbConnected: false,
		}, false
	}
	return &Result{
		DbConnected: true,
	}, true
}
