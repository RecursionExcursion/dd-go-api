package app

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/betbot"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

var fsdStringCompressor = lib.GzipCompressor[betbot.FirstShotData](
	lib.Codec[string]{
		Encode: func(b []byte) (string, error) {
			return lib.BytesToBase64(b), nil
		},
		Decode: func(s string) ([]byte, error) {
			return lib.Base64ToBytes(s)
		},
	},
)

var HandleBBGet api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	_, dataRepo := BetBotRepository()

	timer := lib.StartTimer()

	lib.Log("Querying DB for betbot data", 5)
	compressedData, err := dataRepo.FindTById(dataId)
	if err != nil {
		log.Println(err)
		api.Response.ServerError(w, "")
		return
	}

	lib.Log("Decompressing Data", 5)
	decompressedDbData, err := fsdStringCompressor.Decompress(compressedData.Data)
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	lib.Log("Compiling stats", 5)
	packagedData, err := betbot.NewStatCalculator(decompressedDbData).CalculateAndPackage()
	if err != nil {
		lib.LogError("", err)
		api.Response.ServerError(w)
		return
	}

	lib.Log("Gzipping payload", 5)

	api.Response.Gzip(w, 200,
		struct {
			Meta int64
			Data []betbot.PackagedPlayer
		}{
			Meta: decompressedDbData.Created,
			Data: packagedData,
		},
	)

	timer.End()
}

/* Atomic bool for tracking whether or not the validation process is on going */
var isWorking atomic.Bool

var handleRevalidationPolling api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	api.Response.Ok(w, isWorking.Load())
}

var HandleGetBBRevalidation api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {
	if isWorking.Load() {
		api.Response.Ok(w, "Revalidation in progress")
		return
	}

	isWorking.Store(true)

	/* Long running task, done in the bg, tacked by the isWorking atomic Bool */
	go func() {
		timer := lib.StartTimer()
		defer func() {
			isWorking.Store(false)
			timer.End()
		}()

		//collect data
		lib.Log("Collecting Data", 5)
		fsd, err := betbot.CollectData()
		if err != nil {
			log.Printf("Error while collecting data: %v", err)
			return
		}

		//compress data
		lib.Log("Compressing Data", 5)
		compressedData, err := fsdStringCompressor.Compress(fsd)
		if err != nil {
			log.Printf("Error while compressing data: %v", err)
			return
		}
		compressed := betbot.CompressedFsData{
			Id:      dataId,
			Created: fsd.Created,
			Data:    compressedData,
		}

		_, dataRepo := BetBotRepository()

		//Wipe old data
		lib.Log("Wiping stale Data", 5)
		ok, err := dataRepo.DeleteById(dataId)
		if err != nil || !ok {
			log.Println("Error while wiping data")
			lib.Log(err.Error(), -1)
			return
		}

		//save data
		lib.Log("Saving New Data", 5)
		_, err = dataRepo.SaveT(compressed)
		if err != nil {
			log.Println("Error while saving data")
			lib.Log(err.Error(), -1)
			return
		}

	}()

	api.Response.Ok(w, "Revalidation started")
}

/* Collect, Compute, Send (No state is saved) */
var HandleBBValidateAndZip = func(w http.ResponseWriter, r *http.Request) {
	//collect data
	fsd, err := betbot.CollectData()
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	betbot.FindGameInFsd(fsd, strconv.Itoa(401705613))

	// Compile stats
	packagedData, err := betbot.NewStatCalculator(fsd).CalculateAndPackage()
	if err != nil {
		lib.LogError("", err)
		api.Response.ServerError(w)
		return
	}

	// Zip and return
	lib.Log("Gzipping payload", 5)
	api.Response.Gzip(w, 200,
		struct {
			Meta string
			Data []betbot.PackagedPlayer
		}{
			Meta: strconv.FormatInt(time.Now().UnixMilli(), 10),
			Data: packagedData,
		},
	)

}

// TODO impl pw hashing
var HandleUserLogin api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	type LoginPayload struct {
		Username string
		Password string
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		api.Response.ServerError(w, "Failed to read body")
		return
	}
	defer r.Body.Close()

	pl, err := lib.Map[LoginPayload](bodyBytes)
	if err != nil {
		api.Response.ServerError(w, "Failed to read body")
		return
	}

	/* Validate body */
	if pl.Username == "" || pl.Password == "" {
		api.Response.BadRequest(w)
		return
	}

	userRepo, _ := BetBotRepository()

	usr, err := userRepo.FindFirstT()
	if err != nil {
		api.Response.NotFound(w)
		return
	}

	if usr.Username != pl.Username || usr.Password != pl.Password {
		api.Response.Unauthorized(w)
		return
	}

	claims := map[string]any{
		"sub": usr.Username,
	}

	jwt, err := lib.CreateJWT(claims, time.Hour*48, lib.EnvGetOrPanic("BB_JWT_SECRET"))
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	api.Response.Ok(w, jwt)
}
