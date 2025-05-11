package main

import "net/http"

func (a *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "ok",
		"env":    a.config.env,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		a.internalServerError(w, r, err)
	}
}
