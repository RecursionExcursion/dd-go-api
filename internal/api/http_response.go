package api

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
)

type response = func(http.ResponseWriter, ...any)
type customResponse = func(http.ResponseWriter, int, ...any)

type ApiResponses struct {
	Ok              response
	ServerError     response
	NotFound        response
	Unauthorized    response
	Forbidden       response
	TooManyRequests response
	BadRequest      response
	Send            customResponse
	Gzip            func(w http.ResponseWriter, status int, data ...any)
}

var Response = ApiResponses{
	/* 100 */

	/* 200 */
	Ok: func(w http.ResponseWriter, data ...any) {
		send(w, 200, data)
	},

	/* 300 */

	/* 400 */
	BadRequest: func(w http.ResponseWriter, data ...any) {
		send(w, http.StatusBadRequest, data)
	},

	Unauthorized: func(w http.ResponseWriter, data ...any) {
		send(w, http.StatusUnauthorized, data)
	},

	Forbidden: func(w http.ResponseWriter, data ...any) {
		send(w, http.StatusForbidden, data)
	},

	NotFound: func(w http.ResponseWriter, data ...any) {
		send(w, http.StatusNotFound, data)
	},

	TooManyRequests: func(w http.ResponseWriter, data ...any) {
		send(w, http.StatusTooManyRequests, data)
	},

	/* 500 */
	ServerError: func(w http.ResponseWriter, data ...any) {
		send(w, http.StatusInternalServerError, data)
	},

	/* Misc */
	Send: func(w http.ResponseWriter, status int, data ...any) {
		send(w, status, data)
	},

	Gzip: func(w http.ResponseWriter, status int, data ...any) {
		zip(w, status, data...)
	},
}

func send(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	encodeToJson(data, w)
}

func encodeToJson(data any, w http.ResponseWriter) {
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func zip(w http.ResponseWriter, status int, data ...any) {
	// Set headers
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/plain")

	gz := gzip.NewWriter(w)
	defer gz.Close()

	//json -> gz -> res
	w.WriteHeader(status)
	err := json.NewEncoder(gz).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
