package app

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var atlasUri = lib.EnvGet("ATLAS_URI")

type queryFn[T any] func(c *mongo.Collection) (T, error)

type mongoConnection struct {
	Db         string
	Collection string
}

/* mongoQuery creates, manages, and closes connection while executing
 * the custom query fn passed in
 */
func mongoQuery[T any](mc mongoConnection, query queryFn[T]) (T, error) {
	client, err := connectMongoClient()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	db := client.Database(mc.Db)
	c := db.Collection(mc.Collection)
	return query(c)

}

/* BsonToJson handles converting MongoDb results to json
 * SR- is for SingleResult responses
 */
func bsontoJson() struct {
	SR func(sr *mongo.SingleResult) ([]byte, error)
} {
	return struct {
		SR func(sr *mongo.SingleResult) ([]byte, error)
	}{
		SR: func(sr *mongo.SingleResult) ([]byte, error) {
			var res bson.M
			err := sr.Decode(&res)
			if err == mongo.ErrNoDocuments {
				return nil, errors.New("no document was found with the key")
			}
			if err != nil {
				return nil, err
			}
			jsonData, err := json.MarshalIndent(res, "", "    ")
			if err != nil {
				return nil, err
			}
			return jsonData, nil
		},
	}
}

func connectMongoClient() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(atlasUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}
