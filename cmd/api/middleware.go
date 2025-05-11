package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"ws_practice_1/internal/store"

	"github.com/golang-jwt/jwt/v5"
)

func (app *application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization cookie is missing"))
			return
		}

		token := cookie.Value

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		fmt.Println(user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) ParamAuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := r.URL.Query().Get("token")

		if cookie == "" {
			app.unauthorizedErrorResponse(w, r, fmt.Errorf("authorization cookie is missing"))
			return
		}

		token := cookie

		jwtToken, err := app.authenticator.ValidateToken(token)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		claims, _ := jwtToken.Claims.(jwt.MapClaims)

		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()

		user, err := app.store.Users.GetByID(ctx, userID)
		if err != nil {
			app.unauthorizedErrorResponse(w, r, err)
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		log.Println("CTX:", user.ID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) verifyTokenHandler(w http.ResponseWriter, r *http.Request) {
	app.jsonResponse(w, http.StatusOK, map[string]any{
		"authenticated": true,
	})
}

func (app *application) meHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtx).(*store.User)
	if !ok || user == nil {
		app.unauthorizedErrorResponse(w, r, fmt.Errorf("no user in context"))
		return
	}

	app.jsonResponse(w, http.StatusOK, user)
}
