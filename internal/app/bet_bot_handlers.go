package app

import (
	"net/http"
	"time"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/betbot"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

const dataId = "data"

func HandleGetBetBot(w http.ResponseWriter, r *http.Request) {
	compressedData, err := BetBotRepository().dataRepo.findTById(dataId)
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	data, err := lib.GzipCompressor[betbot.FirstShotData]().Decompress(compressedData.Data)
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	api.Response.Gzip(w, data)
}

func HandleBetBotRevalidation(w http.ResponseWriter, r *http.Request) {

	//collect data
	fsd, err := betbot.CollectData()
	if err != nil {
		api.Response.ServerError(w)
		return
	}

	//compress data
	compressedData, err := lib.GzipCompressor[betbot.FirstShotData]().Compress(fsd)
	if err != nil {
		api.Response.ServerError(w)
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
		api.Response.ServerError(w)
		return
	}

	if ok {
		api.Response.Ok(w)
	} else {
		api.Response.ServerError(w)
	}
}
