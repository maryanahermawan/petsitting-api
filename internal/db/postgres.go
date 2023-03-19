package db

import (
	"context"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

func GetDbConnectionFactory(url string) *pgxpool.Pool {
	dbpool, err := pgxpool.Connect(context.Background(), url)
	if err != nil {
		log.Errorf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return dbpool
}
