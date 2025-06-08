package cfbr

import (
	"github.com/RecursionExcursion/cfbr-core-go/cfbrcore"
	"github.com/RecursionExcursion/go-toolkit/core"
	"github.com/RecursionExcursion/go-toolkit/mongo"
)

type CfbrRepo = mongo.MongoConnection[cfbrcore.SerializeableCompressedSeason]

func CfbrRepository() CfbrRepo {
	dbName := core.EnvGetOrPanic("DB_NAME_CFBR")
	return CfbrRepo{
		Db:         dbName,
		Collection: "seasons",
	}
}
