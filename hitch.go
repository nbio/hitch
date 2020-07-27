package hitch

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Hitch ties httprouter, context, and middleware up in a bow.
type Hitch struct {
	Router     *httprouter.Router
	middleware []func(http.Handler) http.Handler
}

// New initializes a new Hitch.
func New() *Hitch {
	r := httprouter.New()
	r.HandleMethodNotAllowed = false // may cause problems otherwise
	return &Hitch{
		Router: r,
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

// Next registers an http.Handler as a fallback/not-found handler.
func (h *Hitch) Next(handler http.Handler) {
	h.Router.NotFound = handler
}

// Handle registers a handler for the given method and path.
func (h *Hitch) Handle(method, path string, handler http.Handler, middleware ...func(http.Handler) http.Handler) {
	for i := len(middleware) - 1; i >= 0; i-- {
		handler = middleware[i](handler)
	}
	h.Router.Handler(method, path, handler)
}

// HandleFunc registers a func handler for the given method and path.
func (h *Hitch) HandleFunc(method, path string, handler func(http.ResponseWriter, *http.Request), middleware ...func(http.Handler) http.Handler) {
	h.Handle(method, path, http.HandlerFunc(handler), middleware...)
}

// Get registers a GET handler for the given path.
func (h *Hitch) Get(path string, handler http.Handler, middleware ...func(http.Handler) http.Handler) {
	h.Handle("GET", path, handler, middleware...)
}

// Put registers a PUT handler for the given path.
func (h *Hitch) Put(path string, handler http.Handler, middleware ...func(http.Handler) http.Handler) {
	h.Handle("PUT", path, handler, middleware...)
}

// Post registers a POST handler for the given path.
func (h *Hitch) Post(path string, handler http.Handler, middleware ...func(http.Handler) http.Handler) {
	h.Handle("POST", path, handler, middleware...)
}

// Patch registers a PATCH handler for the given path.
func (h *Hitch) Patch(path string, handler http.Handler, middleware ...func(http.Handler) http.Handler) {
	h.Handle("PATCH", path, handler, middleware...)
}

// Delete registers a DELETE handler for the given path.
func (h *Hitch) Delete(path string, handler http.Handler, middleware ...func(http.Handler) http.Handler) {
	h.Handle("DELETE", path, handler, middleware...)
}

// Options registers a OPTIONS handler for the given path.
func (h *Hitch) Options(path string, handler http.Handler, middleware ...func(http.Handler) http.Handler) {
	h.Handle("OPTIONS", path, handler, middleware...)
}

// Handler returns an http.Handler for the embedded router and middleware.
func (h *Hitch) Handler() http.Handler {
	var handler http.Handler = h.Router
	for i := len(h.middleware) - 1; i >= 0; i-- {
		handler = h.middleware[i](handler)
	}
	return handler
}

// Params returns the httprouter.Params for req.
// This is just a passthrough to httprouter.ParamsFromContext.
func Params(req *http.Request) httprouter.Params {
	return httprouter.ParamsFromContext(req.Context())
}
