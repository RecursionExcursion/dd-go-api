package cfbr

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/RecursionExcursion/cfbr-core-go/cfbrcore"
	"github.com/RecursionExcursion/cfbr-core-go/model"
	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/gouse/gouse"
	"github.com/andybalholm/brotli"
	"go.mongodb.org/mongo-driver/bson"
)

type SeasonCache map[int]model.Season

func (sc *SeasonCache) getSeason(yr int) (model.Season, bool) {
	log.Println("Cache accessed")
	szn, ok := (*sc)[yr]
	if ok {
		log.Println("Cached data found")
	}

	return szn, ok
}

func (sc *SeasonCache) cacheSeason(szn model.Season) {
	(*sc)[szn.Year] = szn
}

var cache = make(SeasonCache)

func CfbrRoutes(mwChain []gouse.Middleware) []gouse.RouteHandler {

	cfbrHttpMethods := gouse.NewPathBuilder("/cfbr")

	var postCfrbRoute = gouse.RouteHandler{
		MethodAndPath: cfbrHttpMethods.Methods().POST,
		Handler:       handleCreateCfbStats,
		Middleware:    mwChain,
	}

	var getCfbrRoute = gouse.RouteHandler{
		MethodAndPath: cfbrHttpMethods.Methods().GET,
		Handler:       handleGetCfbrRankings,
		Middleware:    mwChain,
	}
	var deleteCfbrDataRoute = gouse.RouteHandler{
		MethodAndPath: cfbrHttpMethods.Methods().DELETE,
		Handler:       handleDeleteCfbrData,
		Middleware:    mwChain,
	}

	return []gouse.RouteHandler{
		postCfrbRoute,
		getCfbrRoute,
		deleteCfbrDataRoute,
	}

}

var brotCompressor = core.CustomCompressor[model.Season](
	core.Algorithms{
		Writer: func(w io.Writer) (io.WriteCloser, error) {
			return brotli.NewWriterLevel(w, 11), nil
		},
		Reader: func(r io.Reader) (io.Reader, error) {
			return brotli.NewReader(r), nil
		},
	},
	core.Codec[string]{
		Encode: func(b []byte) (string, error) {
			return core.BytesToBase64(b), nil
		},
		Decode: func(s string) ([]byte, error) {
			return core.Base64ToBytes(s)
		},
	},
)

func handleCreateCfbStats(w http.ResponseWriter, r *http.Request) {
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

	szn, err := cfbrcore.CompileSzn(intYr)
	if err != nil {
		gouse.Response.ServerError(w, "Unable to compile season")
		return
	}
	log.Printf("Season %v created with %v schools, %v schools and %v games\n", szn.Year, len(szn.Schedules), len(szn.Teams), len(szn.Games))

	compressedSeason, err := brotCompressor.Compress(szn)
	if err != nil {
		gouse.Response.ServerError(w, "Error while compressing season")
		return
	}

	scs := cfbrcore.SerializeableCompressedSeason{
		Id:               strconv.Itoa(szn.Year),
		Year:             szn.Year,
		CreatedAt:        int(time.Now().UnixMilli()),
		CompressedSeason: compressedSeason,
	}

	cfbrRepo := CfbrRepository()
	ok, err = cfbrRepo.UpsertT(scs, bson.M{"id": scs.Id})
	if !ok || err != nil {
		gouse.Response.ServerError(w, "Error saving data to db")
		return
	}

	log.Println("Season saved")

	// return szn, nil
	gouse.Response.Created(w)
}

func handleGetCfbrRankings(w http.ResponseWriter, r *http.Request) {
	yr, ok, msg := extractYtQueryParam(r)
	if !ok {
		gouse.Response.BadRequest(w, msg)
		return
	}

	intYr, err := strconv.Atoi(yr)
	if err != nil {
		gouse.Response.BadRequest(w, "Year arg cannot be casted to an int")
		return
	}

	/* TODO
	* Only save (cache) curr year and possibly recent years data is about 143kb a szn (compressed)
	 */
	szn, ok := cache.getSeason(intYr)
	if !ok {
		//TODO sanitize inputs year < now.year, div must be valid, etc
		cfbrRepo := CfbrRepository()
		dbSzn, err := cfbrRepo.FindTById(yr)
		if err != nil {
			gouse.Response.NotFound(w, fmt.Sprintf("Season %v not found", yr))
			return
		}
		szn, err = brotCompressor.Decompress(dbSzn.CompressedSeason)
		if err != nil {
			gouse.Response.ServerError(w, "Error during season decompression")
		}
		log.Printf("Season %v found with %v schedules, %v schools and %v games\n", szn.Year, len(szn.Schedules), len(szn.Teams), len(szn.Games))
		cache.cacheSeason(szn)
	}
	tms, gms, err := cfbrcore.MapToRanker(szn)
	if err != nil {
		log.Panic(err)
	}

	s := cfbrcore.RankSeasonProto(tms, gms)

	gouse.Response.Ok(w, []any{s, szn})
}

var handleDeleteCfbrData = func(w http.ResponseWriter, r *http.Request) {
	yr, ok, msg := extractYtQueryParam(r)
	if !ok {
		gouse.Response.BadRequest(w, msg)
		return
	}

	cfbrRepo := CfbrRepository()
	ok, err := cfbrRepo.DeleteById(yr)
	if !ok || err != nil {
		log.Println(err)
		gouse.Response.ServerError(w, "Data could not be deleted")
		return
	}
	gouse.Response.Ok(w, fmt.Sprintf("Season %v DELETED", yr))
}

/* Returns (param, success, error msg) */
func extractYtQueryParam(r *http.Request) (string, bool, string) {
	yr := r.URL.Query().Get("yr")
	if yr == "" {
		return "", false, "Invalid query params"
	}
	return yr, true, ""
}
