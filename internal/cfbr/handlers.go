package cfbr

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/cfbr/core"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
	"go.mongodb.org/mongo-driver/bson"
)

func CfbrRoutes(mwChain []api.Middleware) []api.RouteHandler {

	cfbrHttpMethods := api.HttpMethodGenerator("/cfbr")

	var postCfrbRoute = api.RouteHandler{
		MethodAndPath: cfbrHttpMethods().POST,
		Handler:       handleCreateCfbStats,
		Middleware:    mwChain,
	}

	var getCfbrRoute = api.RouteHandler{
		MethodAndPath: cfbrHttpMethods().GET,
		Handler:       handleGetCfbrRankings,
		Middleware:    mwChain,
	}
	var deleteCfbrDataRoute = api.RouteHandler{
		MethodAndPath: cfbrHttpMethods().DELETE,
		Handler:       handleDeleteCfbrData,
		Middleware:    mwChain,
	}

	return []api.RouteHandler{
		postCfrbRoute,
		getCfbrRoute,
		deleteCfbrDataRoute,
	}

}

var brotCompressor = lib.CustomCompressor[core.Season](
	lib.Algorithms{
		Writer: func(w io.Writer) (io.WriteCloser, error) {
			return brotli.NewWriterLevel(w, 11), nil
		},
		Reader: func(r io.Reader) (io.Reader, error) {
			return brotli.NewReader(r), nil
		},
	},
	lib.Codec[string]{
		Encode: func(b []byte) (string, error) {
			return lib.BytesToBase64(b), nil
		},
		Decode: func(s string) ([]byte, error) {
			return lib.Base64ToBytes(s)
		},
	},
)

func handleCreateCfbStats(w http.ResponseWriter, r *http.Request) {
	yr, ok, msg := extractYtQueryParam(r)
	if !ok {
		api.Response.BadRequest(w, msg)
		return
	}

	intYr, err := strconv.Atoi(yr)
	if err != nil {
		api.Response.BadRequest(w, "Invalid query params")
		return
	}

	szn, err := core.CompileSeason(intYr)
	if err != nil {
		api.Response.ServerError(w, "Unable to compile season")
		return
	}
	log.Printf("Season %v created with %v schools, %v schools and %v games\n", szn.Year, len(szn.Schedules), len(szn.Teams), len(szn.Games))

	compressedSeason, err := brotCompressor.Compress(szn)
	if err != nil {
		api.Response.ServerError(w, "Error while compressing season")
		return
	}

	scs := core.SerializeableCompressedSeason{
		Id:               strconv.Itoa(szn.Year),
		Year:             szn.Year,
		CreatedAt:        int(time.Now().UnixMilli()),
		CompressedSeason: compressedSeason,
	}

	cfbrRepo := CfbrRepository()
	ok, err = cfbrRepo.UpsertT(scs, bson.M{"id": scs.Id})
	if !ok || err != nil {
		api.Response.ServerError(w, "Error saving data to db")
		return
	}

	log.Println("Season saved")

	// return szn, nil
	api.Response.Created(w)
}

func handleGetCfbrRankings(w http.ResponseWriter, r *http.Request) {
	yr, ok, msg := extractYtQueryParam(r)
	if !ok {
		api.Response.BadRequest(w, msg)
		return
	}

	/* TODO
	* Only save (cache) curr year and possibly recent years data is about 143kb a szn (compressed)
	 */
	//TODO sanitize inputs year < now.year, div must be valid, etc
	cfbrRepo := CfbrRepository()
	dbSzn, err := cfbrRepo.FindTById(yr)
	if err != nil {
		api.Response.NotFound(w, fmt.Sprintf("Season %v not found", yr))
		return
	}
	szn, err := brotCompressor.Decompress(dbSzn.CompressedSeason)
	if err != nil {
		api.Response.ServerError(w, "Error during season decompression")
	}
	log.Printf("Season %v found with %v schools, %v schools and %v games\n", szn.Year, len(szn.Schedules), len(szn.Teams), len(szn.Games))

	//TODO compute weights
	// cs, err := core.ComputeSeason(szn)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println("Computation complete")

	tms, gms, err := core.MapToRanker(szn)

	api.Response.Ok(w, []any{tms, gms}, nil)
}

var handleDeleteCfbrData = func(w http.ResponseWriter, r *http.Request) {
	yr, ok, msg := extractYtQueryParam(r)
	if !ok {
		api.Response.BadRequest(w, msg)
		return
	}

	cfbrRepo := CfbrRepository()
	ok, err := cfbrRepo.DeleteById(yr)
	if !ok || err != nil {
		log.Println(err)
		api.Response.ServerError(w, "Data could not be deleted")
		return
	}
	api.Response.Ok(w, fmt.Sprintf("Season %v DELETED", yr))
}

/* Returns (param, success, error msg) */
func extractYtQueryParam(r *http.Request) (string, bool, string) {
	yr := r.URL.Query().Get("yr")
	if yr == "" {
		return "", false, "Invalid query params"
	}
	return yr, true, ""
}
