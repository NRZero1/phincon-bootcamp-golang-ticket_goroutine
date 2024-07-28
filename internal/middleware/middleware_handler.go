package middleware

import "net/http"

type Middleware func(http.Handler) http.Handler

func CreateStack(middlewares ...Middleware) Middleware {
	return func(nextMiddleware http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i -- {
			nextCalledMiddleWare := middlewares[i]
			nextMiddleware = nextCalledMiddleWare(nextMiddleware)
		}

		return nextMiddleware
	}
}