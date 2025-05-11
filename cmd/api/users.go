package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"ws_practice_1/internal/store"

	"github.com/go-chi/chi/v5"
)

type userKey string

const userCtx userKey = "user"

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := app.store.Users.GetByID(ctx, int64(userIDInt))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var updates map[string]interface{}

	if err := readJSON(w, r, &updates); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	err = app.store.Users.Update(ctx, int64(userIDInt), updates)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, "User updated successfully"); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	err = app.store.Users.Delete(ctx, int64(userIDInt))
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err = app.jsonResponse(w, http.StatusOK, "user deleted successfully"); err != nil {
		app.internalServerError(w, r, err)
	}
}

func (app *application) getUserStatsHandler(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(userCtx).(*store.User)
	if !ok || user == nil {
		app.unauthorizedErrorResponse(w, r, fmt.Errorf("no user in context"))
		return
	}

	ctx := r.Context()

	matchesPlayed, err := app.store.Matches.GetMatchesPlayedByUser(ctx, user.ID)
	if err != nil {
		app.internalServerError(w, r, fmt.Errorf("error fetching matches played"))
		return
	}

	matchesWon, err := app.store.Matches.GetMatchesWonByUser(ctx, user.ID)
	if err != nil {
		app.internalServerError(w, r, fmt.Errorf("error fetching matches won"))
		return
	}

	stats := map[string]int{
		"matchesPlayed": matchesPlayed,
		"matchesWon":    matchesWon,
	}

	response := map[string]interface{}{
		"stats": stats,
	}

	if matchesPlayed == 0 && matchesWon == 0 {
		response["message"] = "No matches played yet"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.internalServerError(w, r, fmt.Errorf("error encoding response"))
	}
}

func (app *application) getTotalMatchesPlayedHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	totalMatchesPlayed, err := app.store.Matches.GetTotalMatchesPlayed(ctx)
	if err != nil {
		app.internalServerError(w, r, fmt.Errorf("error fetching total matches played"))
		return
	}

	response := map[string]interface{}{
		"totalMatchesPlayed": totalMatchesPlayed,
	}

	if err := app.jsonResponse(w, http.StatusOK, response); err != nil {
		app.internalServerError(w, r, err)
	}
}
