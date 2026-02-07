package yagaw

import (
	"fmt"
	"iter"
	"maps"
	"net/http"
	"regexp"
)

type RequestHandlerPackage struct {
	Handler     RequestHandler
	IndexParams map[int]string
}
type RequestHandlerMap map[HttpRequestMethod]map[string]RequestHandlerPackage
type RequestHandler func(rw http.ResponseWriter, req *http.Request)

type HttpRequestMethod string

const (
	GET     HttpRequestMethod = `GET`
	HEAD    HttpRequestMethod = `HEAD`
	OPTIONS HttpRequestMethod = `OPTIONS`
	TRACE   HttpRequestMethod = `TRACE`
	PUT     HttpRequestMethod = `PUT`
	DELETE  HttpRequestMethod = `DELETE`
	POST    HttpRequestMethod = `POST`
	PATCH   HttpRequestMethod = `PATCH`
	CONNECT HttpRequestMethod = `CONNECT`
)

type Router struct {
	routes RequestHandlerMap
}

func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	debugRequest(rw, req)
	handler, err := r.findReqHandler(req)
	if err != nil {
		Log.FatalError(err)
	}
	handler(rw, req)
}

func (r *Router) findReqHandler(req *http.Request) (RequestHandler, error) {
	_, methodFound := r.routes[HttpRequestMethod(req.Method)]
	if !methodFound {
		return routeNotFoundHandler, nil
	}

	handlerPackage, routeFound := r.routes[HttpRequestMethod(req.Method)][req.URL.Path]

	if routeFound {
		return handlerPackage.Handler, nil
	}

	if !routeFound {
		path, matchFound := matchRoutePattern(maps.Keys(r.routes[HttpRequestMethod(req.Method)]), req.URL.Path)
		if matchFound {
			return r.routes[HttpRequestMethod(req.Method)][path].Handler, nil
		}
	}

	return routeNotFoundHandler, nil
}

func matchRoutePattern(keysIter iter.Seq[string], path string) (string, bool) {
	for k := range keysIter {
		re := regexp.MustCompile(fmt.Sprintf("(?i)%s", k))
		record := re.FindString(path)
		if len(record) != 0 {
			return k, true
		}
	}
	return "", false
}

func (r *Router) RegisterRoute(method HttpRequestMethod, path string, handler RequestHandler) {
	if r.routes[method] == nil {
		r.routes[method] = make(map[string]RequestHandlerPackage)
	}

	re := regexp.MustCompile(`(?i)({[a-z0-9-_]+})`)
	// If no pattern is found save the route as is
	if re.FindStringIndex(path) == nil {
		r.routes[method][path] = RequestHandlerPackage{Handler: handler}
		return
	}

	// Substitute parameters with the matching regex for future matching
	// Anchor the pattern to match the entire request path
	newPath := "^" + re.ReplaceAllString(path, `([a-z0-9-_]+)`) + "$"
	r.routes[method][newPath] = RequestHandlerPackage{Handler: handler}
}

func (r *Router) RegisteredRoutes() *RequestHandlerMap {
	return &r.routes
}

func NewRouter() *Router {
	return &Router{
		routes: make(RequestHandlerMap),
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
