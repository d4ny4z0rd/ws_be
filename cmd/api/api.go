package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"ws_practice_1/internal/auth"
	"ws_practice_1/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type application struct {
	config        config
	store         store.Storage
	authenticator auth.Authenticator
	ws            wsApp
}

type config struct {
	addr   string
	db     dbConfig
	env    string
	apiURL string
	auth   authConfig
}

type authConfig struct {
	basic basicConfig
	token tokenConfig
}

type basicConfig struct {
	user string
	pass string
}

type tokenConfig struct {
	secret string
	exp    time.Duration
	iss    string
}

type dbConfig struct {
	dbUser       string
	dbPassword   string
	dbHost       string
	dbPort       int
	dbName       string
	dbSSLMode    string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

func (app *application) mount() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		// Update AllowedOrigins to include both HTTP and HTTPS versions if needed
		AllowedOrigins: []string{"http://localhost:5173", "https://stupidcoder.vercel.app", "http://stupidcoder.vercel.app"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Add more headers that might be needed
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Requested-With"},
		ExposedHeaders:   []string{"Set-Cookie", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Credentials"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/health", app.healthCheckHandler)

		r.Route("/users", func(r chi.Router) {

			r.Use(app.AuthTokenMiddleware)

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.AuthTokenMiddleware)
				r.Get("/", app.getUserHandler)
				r.Post("/", app.updateUserHandler)
				r.Delete("/", app.deleteUserHandler)
			})

			r.Get("/stats", app.getUserStatsHandler)
			r.Get("/totalMatchesPlayed", app.getTotalMatchesPlayedHandler)

		})

		r.Route("/ws", func(r chi.Router) {
			r.With(app.AuthTokenMiddleware).Get("/", app.wsHandler)
		})

		r.Route("/authentication", func(r chi.Router) {
			r.Post("/create", app.registerUserHandler)
			r.Post("/token", app.createTokenHandler)
			r.With(app.AuthTokenMiddleware).Get("/me", app.meHandler)
			r.With(app.AuthTokenMiddleware).Get("/verify", app.verifyTokenHandler)
			r.With(app.AuthTokenMiddleware).Post("/logout", app.logoutHandler)
		})
	})

	return r
}

func (app *application) run(mux http.Handler) error {
	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	shutdown := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)

		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		fmt.Printf("Signal caught, %s", s.String())

		shutdown <- srv.Shutdown(ctx)
	}()

	log.Printf("Server has started listening at port %s \n", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdown
	if err != nil {
		return err
	}

	fmt.Printf("Server has stopped, addr: %s, env: %s", app.config.addr, app.config.env)

	return nil
}
