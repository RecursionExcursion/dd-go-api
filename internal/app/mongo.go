package app

import (
	"context"

	"github.com/recursionexcursion/dd-go-api/internal/lib"
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

func connectMongoClient() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(atlasUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return client, nil
}
