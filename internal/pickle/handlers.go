package pickle

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/RecursionExcursion/api-go/api"
	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/go-toolkit/jwt"
	"go.mongodb.org/mongo-driver/bson"
)

func PickleRoutes(mwChain []api.Middleware) []api.RouteHandler {

	picklePlayerEndpoints := api.HttpMethodGenerator("/pickle/player")

	//Player routes
	createPlayer := api.RouteHandler{
		MethodAndPath: picklePlayerEndpoints().POST,
		Handler:       createPlayerHandler,
		Middleware:    mwChain,
	}

	getPlayer := api.RouteHandler{
		MethodAndPath: picklePlayerEndpoints().GET,
		Handler:       getPlayersHandler,
		Middleware:    mwChain,
	}

	// updatePlayer := api.RouteHandler{
	// 	MethodAndPath: picklePlayerEndpoints().PUT,

	// 	Middleware: mwChain,
	// }

	deletePlayer := api.RouteHandler{
		MethodAndPath: picklePlayerEndpoints().DELETE,
		Handler:       deletePlayerHandler,
		Middleware:    mwChain,
	}

	pickleMatchEndpoints := api.HttpMethodGenerator("/pickle/match")

	// //match routes
	createMatch := api.RouteHandler{
		MethodAndPath: pickleMatchEndpoints().POST,
		Handler:       createMatchHandler,
		Middleware:    mwChain,
	}

	getMatch := api.RouteHandler{
		MethodAndPath: pickleMatchEndpoints().GET,
		Handler:       getMatchesHandler,
		Middleware:    mwChain,
	}

	deleteMatch := api.RouteHandler{
		MethodAndPath: pickleMatchEndpoints().DELETE,
		Handler:       deleteMatchHandler,
		Middleware:    mwChain,
	}

	return []api.RouteHandler{
		createPlayer,
		getPlayer,
		// updatePlayer,
		deletePlayer,
		createMatch,
		getMatch,
		deleteMatch,
	}
}

func PickleLoginRoute(mwChain []api.Middleware) []api.RouteHandler {
	pickleAuthEndpoints := api.HttpMethodGenerator("/pickle/auth")

	login := api.RouteHandler{
		MethodAndPath: pickleAuthEndpoints().POST,
		Handler:       loginHandler,
		Middleware:    mwChain,
	}
	return []api.RouteHandler{
		login,
	}

}

var dataId = "pickle"

func getData() PickleData {
	pickleRepo := PickleRepository()
	pd, err := pickleRepo.FindFirstT()
	if err != nil {
		log.Println("Pickle data not found, creating data")
		pickleRepo.SaveT(PickleData{
			ID:      dataId,
			Players: []PicklePlayer{},
			Matches: []PickleMatch{},
		})
		pd, _ = pickleRepo.FindFirstT()
	}
	return pd
}

func saveData(pd PickleData) (bool, error) {
	pickleRepo := PickleRepository()
	return pickleRepo.UpsertT(pd, bson.M{"id": dataId})
}

var createPlayerHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	newPlayerParams, err := readBodyAsJson[struct {
		Name string `json:"name"`
	}](r)
	if err != nil {
		api.Response.BadRequest(w, err.Error())
	}

	if newPlayerParams.Name == "" {
		api.Response.BadRequest(w)
		return
	}

	data := getData()

	data.addPlayer(
		PicklePlayer{
			Id:      generateUID(),
			Name:    newPlayerParams.Name,
			Matches: []string{},
		})

	_, err = saveData(data)
	if err != nil {
		api.Response.ServerError(w, "Failed to save data")
		return
	}

	api.Response.Created(w)
}

var getPlayersHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	data := getData()
	api.Response.Ok(w, data.Players)
}

var deletePlayerHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	deletePlayerParams, err := readBodyAsJson[struct {
		Id string `json:"id"`
	}](r)
	if err != nil {
		api.Response.BadRequest(w, err.Error())
	}

	data := getData()
	ok := data.removePlayer(deletePlayerParams.Id)
	if !ok {
		api.Response.NotFound(w)
	} else {
		saveData(data)
		api.Response.Ok(w)
	}

}

var createMatchHandler = func(w http.ResponseWriter, r *http.Request) {
	p, err := readBodyAsJson[struct {
		Date  int          `json:"date"`
		Score []MatchScore `json:"score"`
	}](r)
	if err != nil {
		api.Response.BadRequest(w, err.Error())
		return
	}

	match := PickleMatch{
		Id:    generateUID(),
		Date:  p.Date,
		Score: p.Score,
	}

	data := getData()
	data.addMatch(match)
	saveData(data)
	api.Response.Ok(w)
}

var getMatchesHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	data := getData()
	api.Response.Ok(w, data.Matches)
}

var deleteMatchHandler = func(w http.ResponseWriter, r *http.Request) {
	p, err := readBodyAsJson[struct {
		Id string `json:"id"`
	}](r)
	if err != nil {
		api.Response.BadRequest(w, err.Error())
		return
	}
	data := getData()
	ok := data.removeMatch(p.Id)
	if !ok {
		api.Response.NotFound(w)
	} else {
		saveData(data)
		api.Response.Ok(w)
	}
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

func generateUID() string {
	nano := time.Now().UnixNano()
	randPart := rand.Intn(1000)
	return strconv.FormatInt(nano+int64(randPart), 36)
}

type LoginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var loginHandler api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	creds, err := readBodyAsJson[LoginPayload](r)
	if err != nil {
		api.Response.BadRequest(w)
		return
	}

	un := core.EnvGetOrPanic("PICKLE_USERNAME")
	pw := core.EnvGetOrPanic("PICKLE_PASSWORD")

	if un != creds.Username || pw != creds.Password {
		api.Response.Unauthorized(w)
		return
	}

	claims := map[string]any{
		"sub": creds.Username,
	}

	jwt, err := jwt.CreateJWT(claims, time.Hour*24, core.EnvGetOrPanic("PICKLE_SECRET"))
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	api.Response.Ok(w, jwt)
}
