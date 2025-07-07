package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"

	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/gouse/gouse"
)

func HashRoutes(mwChain []gouse.Middleware) []gouse.RouteHandler {

	hashEndpoints := gouse.NewPathBuilder("/hash")

	createHash := gouse.RouteHandler{
		MethodAndPath: hashEndpoints.Methods().POST,
		Handler:       postHandler,
		Middleware:    mwChain,
	}

	wakeupHash := gouse.RouteHandler{
		MethodAndPath: hashEndpoints.Methods().GET,
		Handler:       wakeupHandler,
		Middleware:    mwChain,
	}

	return []gouse.RouteHandler{
		createHash,
		wakeupHash,
	}
}

var wakeupHandler gouse.HandlerFn = func(w http.ResponseWriter, _ *http.Request) {
	gouse.Response.Ok(w)
}

type HashRequest struct {
	Message string `json:"message"`
}

type HashResponse struct {
	Hash string `json:"hash"`
}

var postHandler gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	s, err := readBodyAsJson[HashRequest](r)
	if err != nil {
		gouse.Response.BadRequest(w)
		return
	}
	hash := sha256.Sum256([]byte(s.Message))
	resp := HashResponse{Hash: hex.EncodeToString(hash[:])}
	gouse.Response.Created(w, resp)
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
