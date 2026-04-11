package main

import "net/http"

func (app *app) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeErr := writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
	if writeErr != nil {
		return
	}
}

func (app *app) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeErr := writeJSONError(w, http.StatusBadRequest, err.Error())
	if writeErr != nil {
		return
	}
}

func (app *app) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())

	writeErr := writeJSONError(w, http.StatusNotFound, "not found")
	if writeErr != nil {
		return
	}
}
