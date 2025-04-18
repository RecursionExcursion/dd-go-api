package api

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
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
	StreamBytes     func(w http.ResponseWriter, status int, bytes []byte, name string)
	StreamFile      func(w http.ResponseWriter, status int, binPath string, name string)
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

	StreamBytes: func(w http.ResponseWriter, status int, bytes []byte, name string) {
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", name))
		w.WriteHeader(status)
		w.Write(bytes)
	},

	StreamFile: func(w http.ResponseWriter, status int, binPath string, name string) {

		f, _ := os.Open(binPath)
		defer f.Close()

		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\"", name))
		w.WriteHeader(status)

		size, err := io.Copy(w, f)
		if err != nil {
			log.Println("Streaming failed:", err)
		}

		log.Printf("Copied %v bytes", size)

		err = os.RemoveAll(filepath.Dir(binPath))
		if err != nil {
			log.Println("Failed to clean up temp dir:", err)
		}
	},
}

func send(w http.ResponseWriter, status int, data any) {
	switch v := data.(type) {
	case []any:
		if len(v) == 1 {
			data = v[0]
		}
	}

	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func zip(w http.ResponseWriter, status int, data ...any) {
	var payload any
	if len(data) == 1 {
		payload = data[0]
	} else {
		payload = data
	}

	// Set headers
	w.Header().Set("Content-Encoding", "gzip")
	w.Header().Set("Content-Type", "text/plain")

	gz := gzip.NewWriter(w)
	defer gz.Close()

	//json -> gz -> res
	w.WriteHeader(status)
	err := json.NewEncoder(gz).Encode(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
