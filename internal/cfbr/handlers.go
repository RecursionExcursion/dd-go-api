package cfbr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/RecursionExcursion/cfbr-core-go/cfbrcore"
	"github.com/RecursionExcursion/gouse/gouse"
	"github.com/andybalholm/brotli"
)

type ApiSeasonData struct {
	SeasonInfo cfbrcore.SeasonInfo
	GameData   map[string]cfbrcore.GameData
}

type SerializeableCompressedSeason struct {
	Id               string
	Year             int
	CreatedAt        int64
	CompressedSeason []byte
}

//TODO write cache

// type SeasonCache map[int]model.Season

// func (sc *SeasonCache) getSeason(yr int) (model.Season, bool) {
// 	log.Println("Cache accessed")
// 	szn, ok := (*sc)[yr]
// 	if ok {
// 		log.Println("Cached data found")
// 	}

// 	return szn, ok
// }

// func (sc *SeasonCache) cacheSeason(szn model.Season) {
// 	(*sc)[szn.Year] = szn
// }

// var cache = make(SeasonCache)

func CfbrRoutes(mwChain []gouse.Middleware) []gouse.RouteHandler {

	cfbrHttpMethods := gouse.NewPathBuilder("/cfbr")

	var postCfrbRoute = gouse.RouteHandler{
		MethodAndPath: cfbrHttpMethods.Methods().POST,
		Handler:       handleScrapeGameData,
		Middleware:    mwChain,
	}

	var getCfbrRoute = gouse.RouteHandler{
		MethodAndPath: cfbrHttpMethods.Methods().GET,
		Handler:       handleGetSeasonData,
		Middleware:    mwChain,
	}

	// var deleteCfbrDataRoute = gouse.RouteHandler{
	// 	MethodAndPath: cfbrHttpMethods.Methods().DELETE,
	// 	Handler:       handleDeleteCfbrData,
	// 	Middleware:    mwChain,
	// }

	cfbrGameDataPath := gouse.NewPathBuilder("/cfbr/game")

	var getCfbrGameDataRoute = gouse.RouteHandler{
		MethodAndPath: cfbrGameDataPath.Methods().GET,
		Handler:       handleGetSeasonData,
		Middleware:    mwChain,
	}

	return []gouse.RouteHandler{
		postCfrbRoute,
		getCfbrRoute,
		// deleteCfbrDataRoute,
		getCfbrGameDataRoute,
	}

}

var seasonCompressor = struct {
	Compress   func(d ApiSeasonData) ([]byte, error)
	Decompress func([]byte) (ApiSeasonData, error)
}{
	Compress: func(d ApiSeasonData) ([]byte, error) {
		var buf bytes.Buffer

		w := brotli.NewWriterLevel(&buf, 11)

		jsnEncoder := json.NewEncoder(w)
		if err := jsnEncoder.Encode(d); err != nil {
			_ = w.Close()
			return nil, err
		}

		if err := w.Close(); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	},

	Decompress: func(b []byte) (ApiSeasonData, error) {
		var out ApiSeasonData

		r := brotli.NewReader(bytes.NewReader(b))

		dec := json.NewDecoder(r)
		if err := dec.Decode(&out); err != nil {
			return out, err
		}
		return out, nil
	},
}

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
	Level: slog.LevelDebug,
}))

/* Accepts year as arg */
func handleGetSeasonData(w http.ResponseWriter, r *http.Request) {
	yr, ok, msg := extractYrQueryParam(r)
	if !ok {
		gouse.Response.BadRequest(w, msg)
		return
	}

	intYr, err := strconv.Atoi(yr)
	if err != nil {
		gouse.Response.BadRequest(w, "Invalid query params")
		return
	}

	repo, err := CfbrRepository()
	if err != nil {
		log.Println(err)
		gouse.Response.ServerError(w)
		return
	}

	queriedSeason, err := repo.get(yr)
	if err != nil {
		log.Println("Season not found, scraping season data")
		collectedSeason, err := cfbrcore.CollectSeason(intYr, cfbrcore.ScraperOptions{
			Logger: logger,
		})
		if err != nil {
			log.Println(err)
			gouse.Response.ServerError(w, "Error during season collection")
			return
		}

		sznData := ApiSeasonData{
			SeasonInfo: collectedSeason,
			GameData:   map[string]cfbrcore.GameData{},
		}

		err = saveSeason(repo, sznData, intYr)
		if err != nil {
			gouse.Response.ServerError(w)
			return
		}

		gouse.Response.Ok(w, sznData)
		return
	}

	fmt.Println("Season found!")
	decompressedSeason, err := seasonCompressor.Decompress(queriedSeason.CompressedSeason)
	if err != nil {
		log.Println(err)
		gouse.Response.ServerError(w)
		return
	}
	gouse.Response.Ok(w, decompressedSeason.SeasonInfo)
}

func handleScrapeGameData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	yr, ok, msg := extractYrQueryParam(r)
	if !ok {
		gouse.Response.BadRequest(w, msg)
		return
	}

	intYr, err := strconv.Atoi(yr)
	if err != nil {
		gouse.Response.BadRequest(w, "Invalid query params")
		return
	}

	repo, err := CfbrRepository()
	if err != nil {
		log.Println(err)
		gouse.Response.ServerError(w)
		return
	}

	//TODO Clean this shit up!!!!
	var ids []string

	if err := json.NewDecoder(r.Body).Decode(&ids); err != nil {
		gouse.Response.BadRequest(w, " Expected payload of game ids")
		return
	}

	queriedSeason, err := repo.get(yr)
	if err != nil {
		// //TODO cant really happen as the game ids cant be known without seasoninfo being scraped
		gouse.Response.NotFound(w)
		return
	}

	//season found
	decompressedSeason, err := seasonCompressor.Decompress(queriedSeason.CompressedSeason)
	if err != nil {
		log.Println(err)
		gouse.Response.ServerError(w)
		return
	}

	gameData := map[string]cfbrcore.GameData{}
	idsToCollect := []string{}

	fmt.Printf("Games queried %v\n", len(decompressedSeason.GameData))

	for _, id := range ids {
		d, ok := decompressedSeason.GameData[id]
		if !ok {
			idsToCollect = append(idsToCollect, id)
		} else {
			gameData[id] = d
		}
	}

	//Games requested that arent yet in db
	if len(idsToCollect) > 0 {
		gms, err := cfbrcore.CollectGames(idsToCollect, cfbrcore.ScraperOptions{
			Logger: logger,
		})
		if err != nil {
			log.Println(err)
			gouse.Response.ServerError(w)
			return
		}

		//add to return obj and perstence ref
		for _, g := range gms {
			gameData[g.Header.Id] = g
			decompressedSeason.GameData[g.Header.Id] = g
		}

		sznData := ApiSeasonData{
			SeasonInfo: decompressedSeason.SeasonInfo,
			GameData:   decompressedSeason.GameData,
		}

		err = saveSeason(repo, sznData, intYr)
		if err != nil {
			gouse.Response.ServerError(w)
			return
		}
	}

	gouse.Response.Ok(w, gameData)
}

/* Returns (param, success, error msg) */
func extractYrQueryParam(r *http.Request) (string, bool, string) {
	yr := r.URL.Query().Get("yr")
	if yr == "" {
		return "", false, "Invalid query params"
	}
	return yr, true, ""
}

func saveSeason(repo CfbrRepo, sznData ApiSeasonData, intYr int) error {
	log.Println("Inserting season")
	compressed, err := seasonCompressor.Compress(sznData)
	if err != nil {
		log.Println("Error during compressiong")
		log.Println(err)
		return err
	}

	sSzn := SerializeableCompressedSeason{
		Id:               fmt.Sprintf("%v", intYr),
		Year:             intYr,
		CreatedAt:        time.Now().UnixMilli(),
		CompressedSeason: compressed,
	}

	err = repo.insert(sSzn)
	if err != nil {
		log.Println("Error during SQL data insertion")
		log.Println(err)
		return err
	}

	return nil
}
