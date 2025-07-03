package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"

	"github.com/RecursionExcursion/api-go/api"
	"github.com/RecursionExcursion/go-toolkit/core"
)

func HashRoutes(mwChain []api.Middleware) []api.RouteHandler {

	hashEndpoints := api.HttpMethodGenerator("/hash")

	createHash := api.RouteHandler{
		MethodAndPath: hashEndpoints().POST,
		Handler:       postHandler,
		Middleware:    mwChain,
	}

	wakeupHash := api.RouteHandler{
		MethodAndPath: hashEndpoints().GET,
		Handler:       wakeupHandler,
		Middleware:    mwChain,
	}

	return []api.RouteHandler{
		createHash,
		wakeupHash,
	}
}

var wakeupHandler api.HandlerFn = func(w http.ResponseWriter, _ *http.Request) {
	api.Response.Ok(w)
}

type HashRequest struct {
	Message string `json:"message"`
}

type HashResponse struct {
	Hash string `json:"hash"`
}

var postHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	s, err := readBodyAsJson[HashRequest](r)
	if err != nil {
		api.Response.BadRequest(w)
		return
	}
	hash := sha256.Sum256([]byte(s.Message))
	resp := HashResponse{Hash: hex.EncodeToString(hash[:])}
	api.Response.Created(w, resp)
}

func readBodyAsJson[T any](r *http.Request) (T, error) {
	defer r.Body.Close()

	var t T

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return t, errors.New("failed to read body")
	}

	t, err = core.Map[T](bodyBytes)
	if err != nil {
		return t, errors.New("failed to map body")
	}

	return t, nil
}
