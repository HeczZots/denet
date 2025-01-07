package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib" 
	_ "github.com/lib/pq"
)

func InsertMigrations() {
	connRetry := 0
	var migdb *sql.DB
	var err error

	for range time.Tick(time.Second * 5) {
		connRetry++
		if connRetry > 3 {
			log.Fatal("failed to connect to database", "err", err)
		}

		migdb, err = sql.Open("postgres", "postgres://denet:password@db:5432/mydb?sslmode=disable")
		if err != nil {
			continue
		} else {
			break
		}
	}

	driver, err := postgres.WithInstance(migdb, &postgres.Config{
		StatementTimeout: time.Second * 30,
	})
	if err != nil {
		log.Fatal("failed to connect to database for insert migrations", "err", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Could not create migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Could up migrations: %v", err)
	}
	migdb.Close()
	m.Close()
}
