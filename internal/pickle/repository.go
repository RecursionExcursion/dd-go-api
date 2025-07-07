package pickle

import (
	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/go-toolkit/mongo"
)

type PickleRepo = mongo.MongoConnection[PickleData]

func PickleRepository() PickleRepo {
	dbName := core.EnvGetOrPanic("DB_NAME_PICKLE")
	return PickleRepo{
		Db:         dbName,
		Collection: "data",
	}
}
