package core

import "github.com/recursionexcursion/dd-go-api/internal/lib"

const batchSize = 100

var BatchRunner = lib.RunBatchSizeClosure(batchSize)

const BetBotDataId = "data"
