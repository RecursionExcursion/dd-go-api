package app

import (
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/betbot"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

var HandleBBGet api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	lib.Log("Querying DB for betbot data", 5)
	compressedData, err := BetBotRepository().dataRepo.findTById(dataId)
	if err != nil {
		log.Println(err)
		api.Response.ServerError(w, "")
		return
	}

	lib.Log("Decompressing Data", 5)
	decompressedDbData, err := lib.GzipCompressor[betbot.FirstShotData]().Decompress(compressedData.Data)
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	betbot.FindGameInFsd(decompressedDbData, strconv.Itoa(401705610))
	betbot.FindGameInFsd(decompressedDbData, strconv.Itoa(401705611))
	betbot.FindGameInFsd(decompressedDbData, strconv.Itoa(401705612))
	betbot.FindGameInFsd(decompressedDbData, strconv.Itoa(401705613))

	lib.Log("Compiling stats", 5)
	packagedData, err := betbot.NewStatCalculator(decompressedDbData).CalcAndPackage()
	if err != nil {
		lib.LogError("", err)
		api.Response.ServerError(w)
		return
	}

	lib.Log("Gzipping payload", 5)
	api.Response.Gzip(w, 200,
		struct {
			meta string
			data []betbot.PackagedPlayer
		}{
			meta: decompressedDbData.Created,
			data: packagedData,
		},
	)

}

var HandleGetBBRevalidation api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	//collect data
	lib.Log("Collecting Data", 5)
	fsd, err := betbot.CollectData()
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	betbot.FindGameInFsd(fsd, strconv.Itoa(401705613))

	//compress data
	lib.Log("Compressing Data", 5)
	compressedData, err := lib.GzipCompressor[betbot.FirstShotData]().Compress(fsd)
	if err != nil {
		api.Response.ServerError(w, "")
		return
	}
	compressed := betbot.CompressedFsData{
		Id:      dataId,
		Created: time.Now().Format("01-02-2006T15:04:05"),
		Data:    compressedData,
	}

	//Wipe old data
	lib.Log("Wiping stale Data", 5)
	ok, err := BetBotRepository().dataRepo.deleteTById(dataId)
	if err != nil || !ok {
		lib.Log(err.Error(), -1)
		api.Response.ServerError(w, "could not save data")
		return
	}

	//save data
	lib.Log("Saving New Data", 5)
	ok, err = BetBotRepository().dataRepo.saveT(compressed)
	if err != nil {
		lib.Log(err.Error(), -1)
		api.Response.ServerError(w, "could not save data")
		return
	}

	if ok {
		api.Response.Ok(w, "Data revalidated successfully")
	} else {
		api.Response.ServerError(w, "Data could not be revalidated")
	}
}

var HandleBBValidateAndZip = func(w http.ResponseWriter, r *http.Request) {
	//collect data
	fsd, err := betbot.CollectData()
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	log.Println("In handlers")
	betbot.FindGameInFsd(fsd, strconv.Itoa(401705613))

	// Compile stats
	packagedData, err := betbot.NewStatCalculator(fsd).CalcAndPackage()
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

	usr, err := BetBotRepository().userRepo.findFirst()
	if err != nil {
		api.Response.NotFound(w)
		return
	}

	if usr.Username != pl.Username || usr.Password != pl.Password {
		api.Response.Unauthorized(w)
		return
	}

	jwtClaims := make(map[string]any)
	jwtClaims["sub"] = usr.Username

	jwt, err := createJWT(jwtClaims, 48, lib.EnvGet("BB_JWT_SECRET"))
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	api.Response.Ok(w, jwt)
}
