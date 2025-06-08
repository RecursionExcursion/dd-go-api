package betbot

import (
	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/go-toolkit/mongo"
)

type BBUserRepo = mongo.MongoConnection[User]
type BBDataRepo = mongo.MongoConnection[CompressedFsData]

func BetBotRepository() (userRepo BBUserRepo, dataRepo BBDataRepo) {
	dbName := core.EnvGetOrPanic("DB_NAME_BB")

	userRepo = BBUserRepo{
		Db:         dbName,
		Collection: "user",
	}
	dataRepo = BBDataRepo{
		Db:         dbName,
		Collection: "data",
	}

	return userRepo, dataRepo
}
