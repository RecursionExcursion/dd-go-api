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

type CfbrHandler struct {
	Repo *CfbrRepo
}

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

var cache = seasonCache{}

// TODO consdier just using the built in postgress compression
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
// Gets season info
func (h *CfbrHandler) CFBRGet(w http.ResponseWriter, r *http.Request) {
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

	szn, ok := cache.get(intYr)
	if ok {
		log.Println("Using cache")
		gouse.Response.Ok(w, szn.SeasonInfo)
		return
	}

	szn, err = readSeason(h.Repo, yr)
	if err == nil {
		gouse.Response.Ok(w, szn.SeasonInfo)
		return
	}

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

	err = writeSeason(h.Repo, sznData, intYr)
	if err != nil {
		gouse.Response.ServerError(w)
		return
	}

	gouse.Response.Ok(w, sznData)
}

/* gets game data */
func (h *CfbrHandler) CFBRPost(w http.ResponseWriter, r *http.Request) {
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

	ids, err := unmarshalBody[[]string](r)
	if err != nil {
		gouse.Response.BadRequest(w, " Expected payload of game ids")
		return
	}

	/* 	 request has been accepted */
	//check cache
	gameData, missingIds := cache.getGameData(ids)
	if len(missingIds) == 0 {
		gouse.Response.Ok(w, gameData)
		return
	}

	decompressedSeason, err := readSeason(h.Repo, yr)
	if err != nil {
		log.Println(err)
		gouse.Response.ServerError(w)
		return
	}

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

		err = writeSeason(h.Repo, sznData, intYr)
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

func unmarshalBody[T any](r *http.Request) (T, error) {
	var t T
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		return t, err
	}
	return t, nil
}

func readSeason(repo *CfbrRepo, yr string) (ApiSeasonData, error) {
	log.Println("Querying Szn from db")
	queriedSeason, err := repo.get(yr)
	if err != nil {
		return ApiSeasonData{}, err
	}
	szn, err := seasonCompressor.Decompress(queriedSeason.CompressedSeason)
	if err != nil {
		return ApiSeasonData{}, err
	}
	cache.set(szn)
	log.Println("Cache set")
	return szn, nil
}

func writeSeason(repo *CfbrRepo, szn ApiSeasonData, intYr int) error {
	cache.set(szn)
	log.Println("Cache set")

	log.Println("Inserting season into db")
	compressed, err := seasonCompressor.Compress(szn)
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
