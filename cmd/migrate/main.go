package main

import (
	"flag"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"golang-starter/internal/config"
)

func main() {
	direction := flag.String("direction", "up", "up or down")
	flag.Parse()
	if !(*direction == "up" || *direction == "down") {
		log.Fatalf("Invalid argument for direction")
	}

	configuration := config.Init()
	m, err := migrate.New("file://internal/db/migrations", configuration.DbConfig.Url)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if *direction == "up" {
		log.Infof("Running up migration")
		err = m.Up()
	} else if *direction == "down" {
		log.Infof("Running down migration")
		err = m.Down()
	}

	if err == nil {
		version, _, _ := m.Version()
		log.Infof("Migrated to version: %v", version)
	} else if err == migrate.ErrNoChange {
		log.Infof("No change")
	} else {
		log.Fatalf("Failed to run db migrations: %v", err)
	}
}
