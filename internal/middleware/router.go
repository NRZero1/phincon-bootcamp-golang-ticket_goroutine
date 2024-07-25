package middleware

import (
	"fmt"
	"net/http"
)

type Router struct {
	router *http.ServeMux
}

type fn func(http.ResponseWriter, *http.Request)

func NewRouter() *Router {
	return &Router{
		router: http.NewServeMux(),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) AddRoute(method string, url string, handler fn) {
	r.router.HandleFunc(fmt.Sprintf("%s %s", method, url), handler)
}