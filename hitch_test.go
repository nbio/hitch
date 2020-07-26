package hitch

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
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

	if g, e := string(body), "hip-hop"; g != e {
		t.Fatalf("should be == \n \thave: %s\n\twant: %s", g, e)
	}
}

func TestRouteMiddleware(t *testing.T) {
	s := newTestServer(t)
	_, res := s.request("GET", "/route_middleware")
	defer res.Body.Close()
	expectHeaders(t, res)
	body, _ := ioutil.ReadAll(res.Body)

	if g, e := string(body), "middleware1 -> middleware2 -> Hello, world! -> middleware2 -> middleware1"; g != e {
		t.Fatalf("should be == \n \thave: %s\n\twant: %s", g, e)
	}
}

func expectHeaders(t *testing.T, res *http.Response) {
	if g, e := res.Header.Get("Content-Type"), "text/plain"; g != e {
		t.Errorf("should be == \n \thave: %s\n\twant: %s", g, e)
	}
	if g, e := res.Header.Get("X-Awesome"), "awesome"; g != e {
		t.Errorf("should be == \n \thave: %s\n\twant: %s", g, e)
	}
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
	h.Get("/route_middleware", http.HandlerFunc(home), testMiddleware("middleware1"), testMiddleware("middleware2"))

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

func testMiddleware(name string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, name+" -> ")
			next.ServeHTTP(w, req)
			fmt.Fprint(w, " -> "+name)
		})
	}
}
