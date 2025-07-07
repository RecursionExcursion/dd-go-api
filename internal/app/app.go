package app

import "github.com/RecursionExcursion/gouse/gouse"

func App() {
	s := gouse.NewApiServer(":8080", routes())
	s.ListenAndServe()
}
