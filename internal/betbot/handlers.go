package betbot

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/RecursionExcursion/bet-bot-core/bbcore"
	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/go-toolkit/jwt"
	"github.com/RecursionExcursion/gouse/gouse"
)

const BetBotDataId = "data"

func BetbotRoutes(mwChains struct {
	JwtChain []gouse.Middleware
	KeyChain []gouse.Middleware
}) []gouse.RouteHandler {

	bbBase := gouse.NewPathBuilder("/betbot")
	bbPoll := bbBase.Append("poll")
	bbzip := bbBase.Append("zip")
	bbUserLogin := bbBase.Append("user", "login")
	bbPing := bbBase.Append("ping")

	var getBetBotRoute = gouse.RouteHandler{
		MethodAndPath: bbBase.Methods().GET,
		Handler:       handleBBGet,
		Middleware:    mwChains.JwtChain,
	}

	var revalidateBetBotRoute = gouse.RouteHandler{
		MethodAndPath: bbBase.Methods().POST,
		Handler:       handleGetBBRevalidation,
		Middleware:    mwChains.JwtChain,
	}
	var pollBetBotRoute = gouse.RouteHandler{
		MethodAndPath: bbPoll.Methods().GET,
		Handler:       handleRevalidationPolling,
		Middleware:    mwChains.JwtChain,
	}

	var bbRevalidateAndZip = gouse.RouteHandler{
		MethodAndPath: bbzip.Methods().GET,
		Handler:       handleBBValidateAndZip,
		Middleware:    mwChains.JwtChain,
	}

	var loginBBUserRoute = gouse.RouteHandler{
		MethodAndPath: bbUserLogin.Methods().POST,
		Handler:       handleUserLogin,
		Middleware:    mwChains.KeyChain,
	}

	var bbPingRoute = gouse.RouteHandler{
		MethodAndPath: bbPing.Methods().GET,
		Handler: func(w http.ResponseWriter, r *http.Request) {
			gouse.Response.Ok(w)
		},
		Middleware: mwChains.KeyChain,
	}

	return []gouse.RouteHandler{
		getBetBotRoute,
		revalidateBetBotRoute,
		loginBBUserRoute,
		bbRevalidateAndZip,
		bbPingRoute,
		pollBetBotRoute,
	}
}

var fsdStringCompressor = core.GzipCompressor[bbcore.FirstShotData](
	core.Codec[string]{
		Encode: func(b []byte) (string, error) {
			return core.BytesToBase64(b), nil
		},
		Decode: func(s string) ([]byte, error) {
			return core.Base64ToBytes(s)
		},
	},
)

var handleBBGet gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	_, dataRepo := BetBotRepository()

	timer := core.StartTimer()

	log.Println("Querying DB for betbot data", 5)
	compressedData, err := dataRepo.FindTById(BetBotDataId)
	if err != nil {
		log.Println(err)
		gouse.Response.ServerError(w, "")
		return
	}

	log.Println("Decompressing Data", 5)
	decompressedDbData, err := fsdStringCompressor.Decompress(compressedData.Data)
	if err != nil {
		gouse.Response.ServerError(w)
		return
	}

	log.Println("Compiling stats", 5)
	packagedData, err := bbcore.NewStatCalculator(decompressedDbData).CalculateAndPackage()
	if err != nil {
		log.Println("", err)
		gouse.Response.ServerError(w)
		return
	}

	log.Println("Gzipping payload", 5)

	gouse.Response.Gzip(w, 200,
		struct {
			Meta int64
			Data []bbcore.PackagedPlayer
		}{
			Meta: decompressedDbData.Created,
			Data: packagedData,
		},
	)

	timer.End()
}

/* Atomic bool for tracking whether or not the validation process is on going */
var isWorking atomic.Bool

var handleRevalidationPolling gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	gouse.Response.Ok(w, isWorking.Load())
}

var handleGetBBRevalidation gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	if isWorking.Load() {
		gouse.Response.Ok(w, "Revalidation in progress")
		return
	}

	isWorking.Store(true)

	/* Long running task, done in the bg, tacked by the isWorking atomic Bool */
	go func() {
		timer := core.StartTimer()
		defer func() {
			isWorking.Store(false)
			timer.End()
		}()

		//collect data
		log.Println("Collecting Data", 5)
		fsd, err := bbcore.CollectData()
		if err != nil {
			log.Printf("Error while collecting data: %v", err)
			return
		}

		//compress data
		log.Println("Compressing Data", 5)
		compressedData, err := fsdStringCompressor.Compress(fsd)
		if err != nil {
			log.Printf("Error while compressing data: %v", err)
			return
		}
		compressed := CompressedFsData{
			Id:      BetBotDataId,
			Created: fsd.Created,
			Data:    compressedData,
		}

		_, dataRepo := BetBotRepository()

		//Wipe old data
		log.Println("Wiping stale Data", 5)
		ok, err := dataRepo.DeleteById(BetBotDataId)
		if err != nil || !ok {
			log.Println("Error while wiping data")
			log.Println(err.Error(), -1)
			return
		}

		//save data
		log.Println("Saving New Data", 5)
		_, err = dataRepo.SaveT(compressed)
		if err != nil {
			log.Println("Error while saving data")
			log.Println(err.Error(), -1)
			return
		}

	}()

	gouse.Response.Ok(w, "Revalidation started")
}

/* Collect, Compute, Send (No state is saved) */
var handleBBValidateAndZip = func(w http.ResponseWriter, r *http.Request) {
	//collect data
	fsd, err := bbcore.CollectData()
	if err != nil {
		gouse.Response.ServerError(w)
		return
	}

	bbcore.FindGameInFsd(fsd, strconv.Itoa(401705613))

	// Compile stats
	packagedData, err := bbcore.NewStatCalculator(fsd).CalculateAndPackage()
	if err != nil {
		log.Println("", err)
		gouse.Response.ServerError(w)
		return
	}

	// Zip and return
	log.Println("Gzipping payload", 5)
	gouse.Response.Gzip(w, 200,
		struct {
			Meta string
			Data []bbcore.PackagedPlayer
		}{
			Meta: strconv.FormatInt(time.Now().UnixMilli(), 10),
			Data: packagedData,
		},
	)

}

// TODO impl pw hashing
var handleUserLogin gouse.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	type LoginPayload struct {
		Username string
		Password string
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		gouse.Response.ServerError(w, "Failed to read body")
		return
	}
	defer r.Body.Close()

	pl, err := core.Map[LoginPayload](bodyBytes)
	if err != nil {
		gouse.Response.ServerError(w, "Failed to read body")
		return
	}

	/* Validate body */
	if pl.Username == "" || pl.Password == "" {
		gouse.Response.BadRequest(w)
		return
	}

	userRepo, _ := BetBotRepository()

	usr, err := userRepo.FindFirstT()
	if err != nil {
		gouse.Response.NotFound(w)
		return
	}

	if usr.Username != pl.Username || usr.Password != pl.Password {
		gouse.Response.Unauthorized(w)
		return
	}

	claims := map[string]any{
		"sub": usr.Username,
	}

	jwt, err := jwt.CreateJWT(claims, time.Hour*48, core.EnvGetOrPanic("BB_JWT_SECRET"))
	if err != nil {
		gouse.Response.ServerError(w)
		return
	}

	gouse.Response.Ok(w, jwt)
}
