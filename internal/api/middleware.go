package api

import "net/http"

type HandlerFn func(http.ResponseWriter, *http.Request)
type Middleware func(HandlerFn) HandlerFn

/* middlewarePipe passes request through middleware fn from left to right */
func middlewarePipe(mws ...Middleware) Middleware {
	return func(final HandlerFn) HandlerFn {
		for i := len(mws) - 1; i >= 0; i-- {
			final = mws[i](final)
		}
		return final
	}
}
