package hitch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbio/httpcontext"
	"github.com/nbio/st"
)

func TestResponse(t *testing.T) {
	server := createServer()
	defer server.Close()

	print(server.URL + "\n")
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	st.Assert(t, resp.Header.Get("Content-Type"), "text/plain")
}

func createServer() *httptest.Server {
	h := New()
	h.Use(logger, plaintext)
	h.Use(awesome)
	h.HandleFunc("GET", "/", home)
	h.Get("/echo/:phrase", http.HandlerFunc(echo))
	return httptest.NewServer(h)
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		fmt.Printf("%s %s\n", req.Method, req.URL.String())
		next.ServeHTTP(w, req)
	})
}

func plaintext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
	})
}

func awesome(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("X-Awesome", "awesome")
		next.ServeHTTP(w, req)
	})
}

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func echo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, httpcontext.GetString(req, "phrase"))
}
