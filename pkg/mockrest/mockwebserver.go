package mockrest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"time"
)

// A scriptable server which wraps httptest.  Allows us to test the cluster clients by responding back with appropriate schemas

type Server struct {
	testServer *httptest.Server
	handlers   chan http.HandlerFunc
	requests   chan *http.Request
	URL        string
}

func New() *Server {
	return &Server{
		handlers: make(chan http.HandlerFunc),
		requests: make(chan *http.Request),
	}
}

func StartNewWithBody(body string) *Server {
	s := New()
	s.URL = s.Start()
	s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, body)
	})
	return s
}

func StartNewWithStatusCode(status int) *Server {
	s := New()
	s.URL = s.Start()
	s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(status)
	})
	return s
}

func StartNewWithFile(file string) *Server {
	s := New()
	s.URL = s.Start()

	var output string
	b, err := ioutil.ReadFile(file)
	if err != nil {
		output = err.Error()
	} else {
		output = string(b)
	}

	s.Enqueue(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, output)
	})

	return s
}

func (s *Server) Start() string {
	s.testServer = httptest.NewServer(s)
	s.URL = s.testServer.URL
	return s.testServer.URL
}

func (s *Server) Stop() {
	s.testServer.Close()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go func() {
		s.requests <- r
	}()

	select {
	case h := <-s.handlers:
		h.ServeHTTP(w, r)
	default:
		w.WriteHeader(200)
	}
}

func (s *Server) Enqueue(h http.HandlerFunc) {
	go func() {
		s.handlers <- h
	}()
}

func (s *Server) TakeRequest() *http.Request {
	return <-s.requests
}

func (s *Server) TakeRequestWithTimeout(duration time.Duration) *http.Request {
	select {
	case r := <-s.requests:
		return r
	case <-time.After(duration):
		return nil
	}
}
