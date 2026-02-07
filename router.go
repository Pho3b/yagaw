package yagaw

import (
	"fmt"
	"net/http"
)

type ReqHandlerMap map[Method]map[string]ReqHandler
type ReqHandler func(rw http.ResponseWriter, req *http.Request)

type Method string

const (
	GET  Method = "GET"
	POST Method = "POST"
)

type Router struct {
	routes ReqHandlerMap
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	debugRequest(rw, req)
	handler, err := r.findReqHandler(req)
	if err != nil {
		Log.FatalError(err)
	}
	handler(rw, req)
}

func (r *Router) findReqHandler(req *http.Request) (ReqHandler, error) {
	_, methodFound := r.routes[Method(req.Method)]
	if !methodFound {
		return routeNotFoundHandler, nil
	}
	handler, routeFound := r.routes[Method(req.Method)][req.URL.Path]
	if !routeFound {
		return routeNotFoundHandler, nil
	}
	return handler, nil
}

func (r *Router) RegisterRoute(method Method, path string, handler ReqHandler) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]ReqHandler)
	}
	r.routes[method][path] = handler
}

func (r *Router) RegisteredRoutes() *ReqHandlerMap {
	return &r.routes
}

func NewRouter() *Router {
	return &Router{
		routes: make(ReqHandlerMap),
	}
}

func debugRequest(_ http.ResponseWriter, req *http.Request) {
	Log.Debug("Received request:", req.Method, req.URL.Path)
}

func routeNotFoundHandler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/plain")
	rw.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(rw, "404 - Page not found")
}
