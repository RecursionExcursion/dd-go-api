package cfbr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/RecursionExcursion/cfbr-core-go/cfbrcore"
	"github.com/RecursionExcursion/gouse/gouse"
	"github.com/andybalholm/brotli"
)

type ApiSeasonData struct {
	SeasonInfo cfbrcore.SeasonInfo
	GameData   []cfbrcore.GameData
}

type SerializeableCompressedSeason struct {
	Id               string
	Year             int
	CreatedAt        int
	CompressedSeason string
}

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
		Handler:       handleScrapeSeasonInfo,
		Middleware:    mwChain,
	}

	// var deleteCfbrDataRoute = gouse.RouteHandler{
	// 	MethodAndPath: cfbrHttpMethods.Methods().DELETE,
	// 	Handler:       handleDeleteCfbrData,
	// 	Middleware:    mwChain,
	// }

	cfbrScrappingPath := gouse.NewPathBuilder("/cfbr/scrape")

	var getCfbrScrapeRoute = gouse.RouteHandler{
		MethodAndPath: cfbrScrappingPath.Methods().GET,
		Handler:       handleScrapeSeasonInfo,
		Middleware:    mwChain,
	}

	return []gouse.RouteHandler{
		postCfrbRoute,
		getCfbrRoute,
		// deleteCfbrDataRoute,
		getCfbrScrapeRoute,
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
func handleScrapeSeasonInfo(w http.ResponseWriter, r *http.Request) {
	yr, ok, msg := extractYtQueryParam(r)
	if !ok {
		gouse.Response.BadRequest(w, msg)
		return
	}

	intYr, err := strconv.Atoi(yr)
	if err != nil {
		gouse.Response.BadRequest(w, "Invalid query params")
		return
	}

	szn, err := cfbrcore.CollectSeason(intYr, cfbrcore.ScraperOptions{
		Logger: logger,
	})
	if err != nil {
		gouse.Response.ServerError(w, "Unable to compile season")
		return
	}
	// log.Printf("Season %v created with %v schools, %v schools and %v games\n", szn.Year, len(szn.Schedules), len(szn.Teams), len(szn.Games))

	// compressedSeason, err := brotCompressor.Compress(szn)
	// if err != nil {
	// 	gouse.Response.ServerError(w, "Error while compressing season")
	// 	return
	// }

	// scs := cfbrcore.SerializeableCompressedSeason{
	// 	Id:               strconv.Itoa(szn.Year),
	// 	Year:             szn.Year,
	// 	CreatedAt:        int(time.Now().UnixMilli()),
	// 	CompressedSeason: compressedSeason,
	// }

	// cfbrRepo := CfbrRepository()
	// ok, err = cfbrRepo.UpsertT(scs, bson.M{"id": scs.Id})
	// if !ok || err != nil {
	// 	gouse.Response.ServerError(w, "Error saving data to db")
	// 	return
	// }

	// log.Println("Season saved")

	// return szn, nil
	gouse.Response.Ok(w, szn)
}

func handleScrapeGameData(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var p []string

	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		gouse.Response.BadRequest(w, " Expected payload of game ids")
		return
	}

	gms, err := cfbrcore.CollectGames(p, cfbrcore.ScraperOptions{
		Logger: logger})
	if err != nil {
		fmt.Println(err)
		gouse.Response.ServerError(w)
		return
	}

	gouse.Response.Ok(w, gms)
}

// func handleCreateCfbStats(w http.ResponseWriter, r *http.Request) {
// 	yr, ok, msg := extractYtQueryParam(r)
// 	if !ok {
// 		gouse.Response.BadRequest(w, msg)
// 		return
// 	}

// 	intYr, err := strconv.Atoi(yr)
// 	if err != nil {
// 		gouse.Response.BadRequest(w, "Invalid query params")
// 		return
// 	}

// 	szn, err := cfbrcore.CompileSzn(intYr)
// 	if err != nil {
// 		gouse.Response.ServerError(w, "Unable to compile season")
// 		return
// 	}
// 	log.Printf("Season %v created with %v schools, %v schools and %v games\n", szn.Year, len(szn.Schedules), len(szn.Teams), len(szn.Games))

// 	compressedSeason, err := brotCompressor.Compress(szn)
// 	if err != nil {
// 		gouse.Response.ServerError(w, "Error while compressing season")
// 		return
// 	}

// 	scs := SerializeableCompressedSeason{
// 		Id:               strconv.Itoa(szn.Year),
// 		Year:             szn.Year,
// 		CreatedAt:        int(time.Now().UnixMilli()),
// 		CompressedSeason: compressedSeason,
// 	}

// 	cfbrRepo := CfbrRepository()
// 	ok, err = cfbrRepo.UpsertT(scs, bson.M{"id": scs.Id})
// 	if !ok || err != nil {
// 		gouse.Response.ServerError(w, "Error saving data to db")
// 		return
// 	}

// 	log.Println("Season saved")

// 	// return szn, nil
// 	gouse.Response.Created(w)
// }

// func handleGetCfbrRankings(w http.ResponseWriter, r *http.Request) {
// 	yr, ok, msg := extractYtQueryParam(r)
// 	if !ok {
// 		gouse.Response.BadRequest(w, msg)
// 		return
// 	}

// 	intYr, err := strconv.Atoi(yr)
// 	if err != nil {
// 		gouse.Response.BadRequest(w, "Year arg cannot be casted to an int")
// 		return
// 	}

// 	/* TODO
// 	* Only save (cache) curr year and possibly recent years data is about 143kb a szn (compressed)
// 	 */
// 	szn, ok := cache.getSeason(intYr)
// 	if !ok {
// 		//TODO sanitize inputs year < now.year, div must be valid, etc
// 		cfbrRepo := CfbrRepository()
// 		dbSzn, err := cfbrRepo.FindTById(yr)
// 		if err != nil {
// 			gouse.Response.NotFound(w, fmt.Sprintf("Season %v not found", yr))
// 			return
// 		}
// 		szn, err = brotCompressor.Decompress(dbSzn.CompressedSeason)
// 		if err != nil {
// 			gouse.Response.ServerError(w, "Error during season decompression")
// 		}
// 		log.Printf("Season %v found with %v schedules, %v schools and %v games\n", szn.Year, len(szn.Schedules), len(szn.Teams), len(szn.Games))
// 		cache.cacheSeason(szn)
// 	}
// 	tms, gms, err := cfbrcore.MapToRanker(szn)
// 	if err != nil {
// 		log.Panic(err)
// 	}

// 	s := cfbrcore.RankSeasonProto(tms, gms)

// 	gouse.Response.Ok(w, []any{s, szn})
// }

// var handleDeleteCfbrData = func(w http.ResponseWriter, r *http.Request) {
// 	yr, ok, msg := extractYtQueryParam(r)
// 	if !ok {
// 		gouse.Response.BadRequest(w, msg)
// 		return
// 	}

// 	cfbrRepo := CfbrRepository()
// 	ok, err := cfbrRepo.DeleteById(yr)
// 	if !ok || err != nil {
// 		log.Println(err)
// 		gouse.Response.ServerError(w, "Data could not be deleted")
// 		return
// 	}
// 	gouse.Response.Ok(w, fmt.Sprintf("Season %v DELETED", yr))
// }

/* Returns (param, success, error msg) */
func extractYtQueryParam(r *http.Request) (string, bool, string) {
	yr := r.URL.Query().Get("yr")
	if yr == "" {
		return "", false, "Invalid query params"
	}
	return yr, true, ""
}
