package app

import (
	"context"
	"encoding/json"
	"log"

	"github.com/recursionexcursion/dd-go-api/internal/betbot"
	"github.com/recursionexcursion/dd-go-api/internal/lib"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type repo[T any] struct {
	findTById   func(string) (T, error)
	findFirst   func() (T, error)
	saveT       func(T) (bool, error)
	deleteTById func(string) (bool, error)
}

func BetBotRepository() struct {
	userRepo repo[betbot.User]
	dataRepo repo[betbot.CompressedFsData]
} {
	dbName := lib.EnvGet("DB_NAME_BB")

	userConn := mongoConnection{
		Db:         dbName,
		Collection: "user",
	}
	dataConn := mongoConnection{
		Db:         dbName,
		Collection: "data",
	}

	var userRepo = func() repo[betbot.User] {
		return repo[betbot.User]{
			findTById: func(id string) (betbot.User, error) {
				query := bson.D{primitive.E{
					Key:   "id",
					Value: id,
				}}

				return mongoQuery(
					userConn,
					func(c *mongo.Collection) (betbot.User, error) {
						res := c.FindOne(context.TODO(), query)
						jsn, err := bsontoJson().SR(res)
						if err != nil {
							return betbot.User{}, err
						}
						var user betbot.User
						if err := json.Unmarshal(jsn, &user); err != nil {
							log.Println("Error mapping user")
							return betbot.User{}, err
						}
						return user, nil
					})
			},
			saveT: func(user betbot.User) (bool, error) {
				return mongoQuery(userConn, func(c *mongo.Collection) (bool, error) {
					_, err := c.InsertOne(context.TODO(), user)
					if err != nil {
						return false, err
					}
					return true, nil
				})
			},
			findFirst: func() (betbot.User, error) {
				return mongoQuery(
					userConn, func(c *mongo.Collection) (betbot.User, error) {
						// res := c.FindOne(context.TODO())

						jsn, err := bsontoJson().SR(c.FindOne(context.Background(), bson.M{}))
						if err != nil {
							return betbot.User{}, err
						}
						var user betbot.User
						if err := json.Unmarshal(jsn, &user); err != nil {
							log.Println("Error mapping user")
							return betbot.User{}, err
						}

						return user, nil

					},
				)
			},
		}
	}

	var dataRepo = func() repo[betbot.CompressedFsData] {
		return repo[betbot.CompressedFsData]{
			findTById: func(id string) (betbot.CompressedFsData, error) {

				query := bson.D{primitive.E{
					Key:   "id",
					Value: id,
				}}

				return mongoQuery(dataConn, func(c *mongo.Collection) (betbot.CompressedFsData, error) {
					res := c.FindOne(context.TODO(), query)
					// jsn, err := bsontoJson().SR(res)
					// if err != nil {
					// 	return betbot.CompressedFsData{}, err
					// }
					// var fsData betbot.CompressedFsData
					// if err := json.Unmarshal(jsn, &fsData); err != nil {
					// 	log.Println("Error mapping user")
					// 	return betbot.CompressedFsData{}, err
					// }

					var data betbot.CompressedFsData
					err := res.Decode(&data)
					if err == mongo.ErrNoDocuments {
						log.Println("No document found")
						return betbot.CompressedFsData{}, nil // or return error, depending on your app
					}
					if err != nil {
						log.Println("Error decoding from Mongo:", err)
						return betbot.CompressedFsData{}, err
					}
					return data, nil

					// return fsData, nil
				})

			},

			saveT: func(d betbot.CompressedFsData) (bool, error) {
				return mongoQuery(dataConn, func(c *mongo.Collection) (bool, error) {
					_, err := c.InsertOne(context.TODO(), d)
					if err != nil {
						return false, err
					}
					return true, nil
				})
			},
			deleteTById: func(id string) (bool, error) {
				return mongoQuery(dataConn, func(c *mongo.Collection) (bool, error) {
					query := bson.D{primitive.E{
						Key:   "id",
						Value: id,
					}}
					_, err := c.DeleteOne(context.Background(), query)
					if err != nil {
						return false, err
					}
					return true, nil
				})
			},
		}

	}

	return struct {
		userRepo repo[betbot.User]
		dataRepo repo[betbot.CompressedFsData]
	}{
		userRepo: userRepo(),
		dataRepo: dataRepo(),
	}
}
