package app

import (
	"fmt"
	"runtime"

	"github.com/recursionexcursion/dd-go-api/internal/api"
	"github.com/recursionexcursion/dd-go-api/internal/lib"
)

func App() {
	logOsInfo()

	s := api.NewApiServer(":8080")
	s.Init(routes())
	s.ListenAndServe()
}

func logOsInfo() {
	lib.Log("\nOS INFO:", -1)
	lib.Log(fmt.Sprintf("CPUs available: %d", runtime.NumCPU()), -1)
	lib.Log(fmt.Sprintf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0)), -1)
}
