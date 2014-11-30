package hitch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/nbio/st"
)

func TestHome(t *testing.T) {
	s := newTestServer(t)
	_, res := s.request("GET", "/")
	defer res.Body.Close()
	expectHeaders(t, res)
}

func TestEcho(t *testing.T) {
	s := newTestServer(t)
	_, res := s.request("GET", "/api/echo/hip-hop")
	defer res.Body.Close()
	expectHeaders(t, res)
	body, _ := ioutil.ReadAll(res.Body)
	st.Assert(t, string(body), "hip-hop")
}

func expectHeaders(t *testing.T, res *http.Response) {
	st.Expect(t, res.Header.Get("Content-Type"), "text/plain")
	st.Expect(t, res.Header.Get("X-Awesome"), "awesome")
}

// testServer

type testServer struct {
	*httptest.Server
	t *testing.T
}

func (s *testServer) request(method, path string) (*http.Request, *http.Response) {
	req, err := http.NewRequest(method, s.URL+path, nil)
	if err != nil {
		s.t.Fatal(err)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		s.t.Fatal(err)
	}
	return req, res
}

func newTestServer(t *testing.T) *testServer {
	h := New()
	h.Use(logger, plaintext)
	h.UseHandler(http.HandlerFunc(awesome))
	h.HandleFunc("GET", "/", home)
	api := New()
	api.Get("/api/echo/:phrase", http.HandlerFunc(echo))
	h.Next(api.Handler())

	s := &testServer{httptest.NewServer(h.Handler()), t}
	runtime.SetFinalizer(s, func(s *testServer) { s.Server.Close() })
	return s
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
	fmt.Fprint(w, Params(req).ByName("phrase"))
}
