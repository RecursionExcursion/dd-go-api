package pickle

import (
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/go-toolkit/jwt"
	"github.com/RecursionExcursion/gouse/gouse"
	"go.mongodb.org/mongo-driver/bson"
)

func PickleRoutes(mwChain []gouse.Middleware) []gouse.RouteHandler {

	picklePlayerEndpoints := gouse.NewPathBuilder("/pickle/player")

	//Player routes
	createPlayer := gouse.RouteHandler{
		MethodAndPath: picklePlayerEndpoints.Methods().POST,
		Handler:       createPlayerHandler,
		Middleware:    mwChain,
	}

	getPlayer := gouse.RouteHandler{
		MethodAndPath: picklePlayerEndpoints.Methods().GET,
		Handler:       getPlayersHandler,
		Middleware:    mwChain,
	}

	// updatePlayer := api.RouteHandler{
	// 	MethodAndPath: picklePlayerEndpoints().PUT,

	// 	Middleware: mwChain,
	// }

	deletePlayer := gouse.RouteHandler{
		MethodAndPath: picklePlayerEndpoints.Methods().DELETE,
		Handler:       deletePlayerHandler,
		Middleware:    mwChain,
	}

	pickleMatchEndpoints := gouse.NewPathBuilder("/pickle/match")

	// //match routes
	createMatch := gouse.RouteHandler{
		MethodAndPath: pickleMatchEndpoints.Methods().POST,
		Handler:       createMatchHandler,
		Middleware:    mwChain,
	}

	getMatch := gouse.RouteHandler{
		MethodAndPath: pickleMatchEndpoints.Methods().GET,
		Handler:       getMatchesHandler,
		Middleware:    mwChain,
	}

	deleteMatch := gouse.RouteHandler{
		MethodAndPath: pickleMatchEndpoints.Methods().DELETE,
		Handler:       deleteMatchHandler,
		Middleware:    mwChain,
	}

	return []gouse.RouteHandler{
		createPlayer,
		getPlayer,
		// updatePlayer,
		deletePlayer,
		createMatch,
		getMatch,
		deleteMatch,
	}
}

func PickleLoginRoute(mwChain []gouse.Middleware) []gouse.RouteHandler {
	pickleAuthEndpoints := gouse.NewPathBuilder("/pickle/auth")

	login := gouse.RouteHandler{
		MethodAndPath: pickleAuthEndpoints.Methods().POST,
		Handler:       loginHandler,
		Middleware:    mwChain,
	}
	return []gouse.RouteHandler{
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

var createPlayerHandler gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	newPlayerParams, err := gouse.DecodeJSON[struct {
		Name string `json:"name"`
	}](r)
	if err != nil {
		gouse.Response.BadRequest(w, err.Error())
	}

	if newPlayerParams.Name == "" {
		gouse.Response.BadRequest(w)
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
		gouse.Response.ServerError(w, "Failed to save data")
		return
	}

	gouse.Response.Created(w)
}

var getPlayersHandler gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	data := getData()
	gouse.Response.Ok(w, data.Players)
}

var deletePlayerHandler gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	deletePlayerParams, err := gouse.DecodeJSON[struct {
		Id string `json:"id"`
	}](r)
	if err != nil {
		gouse.Response.BadRequest(w, err.Error())
	}

	data := getData()
	ok := data.removePlayer(deletePlayerParams.Id)
	if !ok {
		gouse.Response.NotFound(w)
	} else {
		saveData(data)
		gouse.Response.Ok(w)
	}

}

var createMatchHandler = func(w http.ResponseWriter, r *http.Request) {
	p, err := gouse.DecodeJSON[struct {
		Date  int          `json:"date"`
		Score []MatchScore `json:"score"`
	}](r)
	if err != nil {
		gouse.Response.BadRequest(w, err.Error())
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
	gouse.Response.Ok(w)
}

var getMatchesHandler gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	data := getData()
	gouse.Response.Ok(w, data.Matches)
}

var deleteMatchHandler = func(w http.ResponseWriter, r *http.Request) {
	p, err := gouse.DecodeJSON[struct {
		Id string `json:"id"`
	}](r)
	if err != nil {
		gouse.Response.BadRequest(w, err.Error())
		return
	}
	data := getData()
	ok := data.removeMatch(p.Id)
	if !ok {
		gouse.Response.NotFound(w)
	} else {
		saveData(data)
		gouse.Response.Ok(w)
	}
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

var loginHandler gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	creds, err := gouse.DecodeJSON[LoginPayload](r)
	if err != nil {
		gouse.Response.BadRequest(w)
		return
	}

	un := core.EnvGetOrPanic("PICKLE_USERNAME")
	pw := core.EnvGetOrPanic("PICKLE_PASSWORD")

	if un != creds.Username || pw != creds.Password {
		gouse.Response.Unauthorized(w)
		return
	}

	claims := map[string]any{
		"sub": creds.Username,
	}

	jwt, err := jwt.CreateJWT(claims, time.Hour*24, core.EnvGetOrPanic("PICKLE_SECRET"))
	if err != nil {
		gouse.Response.ServerError(w)
		return
	}

	gouse.Response.Ok(w, jwt)
}
