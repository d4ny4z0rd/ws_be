package main

import (
	"log"
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Internal server error: %s path:%s error:%s \n", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Bad request error: %s path:%s error:%s \n", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Not found error: %s path:%s error:%s \n", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Conflict error: %s path:%s error:%s \n", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Unauthorized error: %s path:%s error:%s \n", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("Forbidden error: %s path:%s error:%s \n", r.Method, r.URL.Path, err.Error())

	writeJSONError(w, http.StatusForbidden, "forbidden")
}
