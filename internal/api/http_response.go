package api

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
)

type ApiResponses struct {
	Ok           func(w http.ResponseWriter)
	OkPayload    func(w http.ResponseWriter, data any)
	Gzip         func(w http.ResponseWriter, data any)
	ServerError  func(w http.ResponseWriter)
	NotFound     func(w http.ResponseWriter)
	Unauthorized func(w http.ResponseWriter)
	Forbidden    func(w http.ResponseWriter)
}

var Response = ApiResponses{
	OkPayload: func(w http.ResponseWriter, data any) {
		w.WriteHeader(http.StatusOK)

		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	},

	Ok: func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusOK)
	},

	Gzip: func(w http.ResponseWriter, data any) {
		// Set headers
		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "text/plain")

		gz := gzip.NewWriter(w)
		defer gz.Close()

		//json -> gz -> res
		err := json.NewEncoder(gz).Encode(data)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	},

	ServerError: func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusInternalServerError)
	},

	NotFound: func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusNotFound)
	},

	Unauthorized: func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusUnauthorized)
	},

	Forbidden: func(w http.ResponseWriter) {
		w.WriteHeader(http.StatusForbidden)
	},
}
