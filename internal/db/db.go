package db

import (
	"context"
	"database/sql"
	"time"
	"ws_practice_1/internal/env"

	_ "github.com/lib/pq"
)

func New(dbUser, dbPassword, dbHost string, dbPort int, dbName, dbSSLMode string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	connectionString := env.GetString("DB_ADDR", "localhost:5432")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)

	duration, err := time.ParseDuration(maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	db.SetMaxIdleConns(maxIdleConns)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	return db, nil
}
