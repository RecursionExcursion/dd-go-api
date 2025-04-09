package betbot

import "github.com/recursionexcursion/dd-go-api/internal/lib"

const batchSize = 100

var BatchRunner = lib.RunBatchSizeClosure(batchSize)
