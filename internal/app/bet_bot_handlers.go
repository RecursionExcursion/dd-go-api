package app

import (
	"net/http"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/betbot"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

const dataId = "data"

var HandleGetBetBot api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	lib.Log("Querying DB for betbot data", 5)

	compressedData, err := BetBotRepository().dataRepo.findTById(dataId)
	if err != nil {
		api.Response.ServerError(w, "")
		return
	}

	lib.Log("Decompressing Data", 5)
	data, err := lib.GzipCompressor[betbot.FirstShotData]().Decompress(compressedData.Data)
	if err != nil {
		api.Response.ServerError(w, "")
		return
	}

	lib.Log("Gzipping payload", 5)
	api.Response.Gzip(w, 200, data)
}

var HandleBetBotRevalidation api.HandlerFn = func(w http.ResponseWriter, r *http.Request) {

	//collect data
	fsd, err := betbot.CollectData()
	if err != nil {
		api.Response.ServerError(w, "")
		return
	}

	//compress data
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

	//save data
	ok, err := BetBotRepository().dataRepo.saveT(compressed)
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
