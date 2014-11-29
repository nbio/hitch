package hitch

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbio/httpcontext"
	"github.com/nbio/st"
)

func TestGet(t *testing.T) {
	server := createServer()
	defer server.Close()

	print(server.URL + "\n")
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	assertHeaders(t, res)
}

func assertHeaders(t *testing.T, res *http.Response) {
	st.Assert(t, res.Header.Get("Content-Type"), "text/plain")
	st.Assert(t, res.Header.Get("X-Awesome"), "awesome")
}

func createServer() *httptest.Server {
	h := New()
	h.Use(logger, plaintext)
	h.UseHandler(http.HandlerFunc(awesome))
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
		next.ServeHTTP(w, req)
	})
}

func awesome(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("X-Awesome", "awesome")
}

func home(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}

func echo(w http.ResponseWriter, req *http.Request) {
	fmt.Fprint(w, httpcontext.GetString(req, "phrase"))
}
