package main

import "net/http"

func (app *app) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	writeJSONError(w, http.StatusInternalServerError, "the server encountered a problem")
}
