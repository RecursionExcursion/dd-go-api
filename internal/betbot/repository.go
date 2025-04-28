package betbot

import (
	"github.com/recursionexcursion/dd-go-api/internal/betbot/core"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

type BBUserRepo = lib.MongoConnection[core.User]
type BBDataRepo = lib.MongoConnection[core.CompressedFsData]

func BetBotRepository() (userRepo BBUserRepo, dataRepo BBDataRepo) {
	dbName := lib.EnvGetOrPanic("DB_NAME_BB")

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
