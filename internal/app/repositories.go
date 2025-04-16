package app

import (
	"github.com/recursionexcursion/dd-go-api/internal/betbot"
	"github.com/recursionexcursion/dd-go-api/internal/cfbr"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

type BBUserRepo = lib.MongoConnection[betbot.User]
type BBDataRepo = lib.MongoConnection[betbot.CompressedFsData]

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

type CfbrRepo = lib.MongoConnection[cfbr.SerializeableCompressedSeason]

func CfbrRepository() CfbrRepo {
	dbName := lib.EnvGetOrPanic("DB_NAME_CFBR")
	return CfbrRepo{
		Db:         dbName,
		Collection: "seasons",
	}
}
