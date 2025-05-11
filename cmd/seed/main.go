package main

import (
	"log"
	"ws_practice_1/internal/db"
	"ws_practice_1/internal/env"
	"ws_practice_1/internal/store"
)

func main() {
	dbConn, err := db.New(
		env.GetString("DB_USER", "admin"),
		env.GetString("DB_PASSWORD", "adminpassword"),
		env.GetString("DB_HOST", "localhost"),
		env.GetInt("DB_PORT", 5432),
		env.GetString("DB_NAME", "ws1"),
		env.GetString("DB_SSL", "disable"),
		env.GetInt("DB_MAX_CONN_OPEN", 30),
		env.GetInt("DB_MAX_IDLE_CONNS", 30),
		env.GetString("DB_MAX_IDLE_TIME", "15m"),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close()

	store := store.NewStorage(dbConn)
	db.Seed(store, dbConn)
}
