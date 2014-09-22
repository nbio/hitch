package hitch

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/nbio/httpcontext"
)

// Hitch ties httprouter, httpcontext, and middleware up in a bow.
type Hitch struct {
	Router *httprouter.Router
	middleware []func(http.Handler) http.Handler
}

// New initializes a new Hitch.
func New() *Hitch {
	return &Hitch{
		Router: httprouter.New(),
	}
}

// Use installs one or more middleware in the Hitch request cycle.
func (h *Hitch) Use(middleware ...func(http.Handler) http.Handler) {
	h.middleware = append(h.middleware, middleware...)
}

// UseHandler registers an http.Handler as a middleware.
func (h *Hitch) UseHandler(handler http.Handler) {
	h.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			handler.ServeHTTP(w, req)
			next.ServeHTTP(w, req)
		})
	})
}

// Handle registers a handler for the given method and path.
func (h *Hitch) Handle(method, path string, handler http.Handler) {
	h.Router.Handle(method, path, wrap(handler))
}

// HandleFunc registers a func handler for the given method and path.
func (h *Hitch) HandleFunc(method, path string, handler func(http.ResponseWriter, *http.Request)) {
	h.Router.Handle(method, path, wrap(http.HandlerFunc(handler)))
}

// Get registers a GET handler for the given path.
func (h *Hitch) Get(path string, handler http.Handler) { h.Handle("GET", path, handler) }

// Put registers a PUT handler for the given path.
func (h *Hitch) Put(path string, handler http.Handler) { h.Handle("PUT", path, handler) }

// Post registers a POST handler for the given path.
func (h *Hitch) Post(path string, handler http.Handler) { h.Handle("POST", path, handler) }

// Patch registers a PATCH handler for the given path.
func (h *Hitch) Patch(path string, handler http.Handler) { h.Handle("PATCH", path, handler) }

// Delete registers a DELETE handler for the given path.
func (h *Hitch) Delete(path string, handler http.Handler) { h.Handle("DELETE", path, handler) }

// Options registers a OPTIONS handler for the given path.
func (h *Hitch) Options(path string, handler http.Handler) { h.Handle("OPTIONS", path, handler) }

// ServeHTTP implements http.Handler.
func (h *Hitch) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	next := http.Handler(h.Router)
	for i := len(h.middleware) - 1; i >= 0; i-- {
		next = h.middleware[i](next)
	}
	next.ServeHTTP(w, req)
}

type key int

const paramsKey key = 1

func wrap(handler http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, params httprouter.Params) {
		httpcontext.Set(req, paramsKey, params)
		handler.ServeHTTP(w, req)
	}
}

// Params returns the httprouter.Params for req.
func Params(req *http.Request) httprouter.Params {
	if value, ok := httpcontext.GetOk(req, paramsKey); ok {
		if params, ok := value.(httprouter.Params); ok {
			return params
		}
	}
	return httprouter.Params{}
}
