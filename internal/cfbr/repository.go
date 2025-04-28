package cfbr

import (
	"github.com/recursionexcursion/dd-go-api/internal/cfbr/core"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

type CfbrRepo = lib.MongoConnection[core.SerializeableCompressedSeason]

func CfbrRepository() CfbrRepo {
	dbName := lib.EnvGetOrPanic("DB_NAME_CFBR")
	return CfbrRepo{
		Db:         dbName,
		Collection: "seasons",
	}
}
