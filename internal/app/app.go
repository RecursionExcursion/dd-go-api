package app

import "github.com/RecursionExcursion/api-go/api"

func App() {
	s := api.NewApiServer(":8080")
	s.Init(routes())
	s.ListenAndServe()
}
