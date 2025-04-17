package app

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/cfbr"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
	"go.mongodb.org/mongo-driver/bson"
)

var gzipCompressor = lib.GzipCompressor[cfbr.SerializableSchoolMap](
	lib.Codec[string]{
		Encode: func(b []byte) (string, error) {
			return lib.BytesToBase64(b), nil
		},
		Decode: func(s string) ([]byte, error) {
			return lib.Base64ToBytes(s)
		},
	},
)

// TODO decide which compression best fits
var brotCompressor = lib.CustomCompressor[cfbr.SerializableSchoolMap](
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

func handleCfbrGet(w http.ResponseWriter, r *http.Request) {
	div := "fbs"
	yr := 2024

	season, err := cfbr.Create(div, uint(yr))
	if err != nil {
		panic(err)
	}

	serSeason := season.Save()

	compressedSeason, err := brotCompressor.Compress(serSeason)
	if err != nil {
		panic(err)
	}

	scs := cfbr.SerializeableCompressedSeason{
		Id:               fmt.Sprintf("%v%v", div, yr),
		Year:             yr,
		CreatedAt:        int(time.Now().UnixMilli()),
		CompressedSeason: compressedSeason,
	}

	cfbrRepo := CfbrRepository()

	cfbrRepo.UpsertT(scs, bson.M{"id": scs.Id})

	api.Response.Ok(w)
}
