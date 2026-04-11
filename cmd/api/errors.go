package main

import "net/http"

func (app *app) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Errorw("internal error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeErr := writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
	if writeErr != nil {
		return
	}
}
