package main

import (
	"log"
	"time"
	"ws_practice_1/internal/auth"
	"ws_practice_1/internal/db"
	"ws_practice_1/internal/env"
	"ws_practice_1/internal/store"

	"github.com/gorilla/websocket"
)

func main() {
	cfg := config{
		addr: env.GetString("PORT", "0.0.0.0:8080"),
		db: dbConfig{
			dbUser:       env.GetString("DB_USER", "admin"),
			dbPassword:   env.GetString("DB_PASSWORD", "adminpassword"),
			dbHost:       env.GetString("DB_HOST", "localhost"),
			dbPort:       env.GetInt("DB_PORT", 5432),
			dbName:       env.GetString("DB_NAME", "ws1"),
			dbSSLMode:    env.GetString("DB_SSL", "disable"),
			maxOpenConns: env.GetInt("DB_MAX_CONN_OPEN", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env:    env.GetString("ENV", "development"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8080"),
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3,
				iss:    "ws1",
			},
		},
	}

	db, err := db.New(
		cfg.db.dbUser,
		cfg.db.dbPassword,
		cfg.db.dbHost,
		cfg.db.dbPort,
		cfg.db.dbName,
		cfg.db.dbSSLMode,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	log.Println("Database connection pool established")

	store := store.NewStorage(db)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		store:         store,
		authenticator: jwtAuthenticator,
	}

	app.ws = wsApp{
		matchMaking: MatchMaking{},
		matches:     make(map[*websocket.Conn]*Match),
		scores:      make(map[*websocket.Conn]int),
		app:         app,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
