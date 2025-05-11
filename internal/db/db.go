package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

func New(dbUser, dbPassword, dbHost string, dbPort int, dbName, dbSSLMode string, maxOpenConns, maxIdleConns int, maxIdleTime string) (*sql.DB, error) {
	connectionString := fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s sslmode=%s",
		dbUser,
		dbPassword,
		dbHost,
		dbPort,
		dbName,
		dbSSLMode,
	)

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
