package app

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/andybalholm/brotli"
	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/cfbr"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
	"go.mongodb.org/mongo-driver/bson"
)

// var gzipCompressor = lib.GzipCompressor[cfbr.CFBRSeason](
// 	lib.Codec[string]{
// 		Encode: func(b []byte) (string, error) {
// 			return lib.BytesToBase64(b), nil
// 		},
// 		Decode: func(s string) ([]byte, error) {
// 			return lib.Base64ToBytes(s)
// 		},
// 	},
// )

// TODO decide which compression best fits
var brotCompressor = lib.CustomCompressor[cfbr.CFBRSeason](
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

	//TODO placeholders
	//TODO sanitize inputs year < now.year, div must be valid, etc
	div := "fbs"
	yr := 20245

	szn := func() cfbr.CFBRSeason {
		cfbrRepo := CfbrRepository()
		queryId := createQueryId(yr, div)

		dbSzn, err := cfbrRepo.FindTById(queryId)
		if err != nil {

			log.Println("Season not found creating new")

			szn, err := cfbr.Create(div, uint(yr))
			if err != nil {
				panic(err)
			}

			/* TODO
			* Only save (cache) curr year and possilby recent years data is about 143kb a szn (compressed)
			 */
			compressedSeason, err := brotCompressor.Compress(szn)
			if err != nil {
				//TODO
				panic(err)
			}

			scs := cfbr.SerializeableCompressedSeason{
				Id:               createQueryId(szn.Year, szn.Division),
				Year:             szn.Year,
				CreatedAt:        int(time.Now().UnixMilli()),
				CompressedSeason: compressedSeason,
			}

			cfbrRepo.UpsertT(scs, bson.M{"id": scs.Id})

			log.Println("Season saved")

			return szn
		} else {
			log.Println("Season found decompressing")
			szn, err := brotCompressor.Decompress(dbSzn.CompressedSeason)
			if err != nil {
				//TODO
				panic(err)
			}
			return szn
		}
	}()

	fmt.Printf("Season %v %v created with %v schools and %v games", szn.Year, szn.Division, len(szn.Schools), len(szn.Games))
	//TODO compute weights

	api.Response.Ok(w)
}

func createQueryId(year int, division string) string {
	return fmt.Sprintf("%v%v", division, year)
}
