package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	DEBUG = true
)

type Server struct {
	Router *Router
}

// constructor for a Server
func NewServer() *Server {
	return &Server{
		Router: NewRouter(),
	}
}

// wrap marshal types for convenience
func marshal(v interface{}) ([]byte, error) {
	if DEBUG {
		return json.MarshalIndent(v, "", "  ")
	} else {
		return json.Marshal(v)
	}
}

func (s *Server) summary(w http.ResponseWriter, r *http.Request) {
	dat, err := marshal(s.Router.Nodes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(dat)

	}
}

func (s *Server) nodes(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/nodes/"):]
	v, ok := s.Router.Nodes[name]
	if ok {
		dat, err := marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write(dat)
		}
	} else {
		http.Error(w, "node not found!", http.StatusNotFound)
	}
}

func (s *Server) types(w http.ResponseWriter, r *http.Request) {
	names := make([]string, 0, len(AvailableNodes))
	for k := range AvailableNodes {
		names = append(names, k)
	}

	dat, err := marshal(names)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.Write(dat)
	}

}

// start a server serving. blocks until err or exit
func (s *Server) Start() error {
	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/":        s.summary,
		"/summary": s.summary,
		"/nodes/":  s.nodes,
		"/types":   s.types,
	}

	for ep, h := range handlers {
		http.HandleFunc(ep, h)
	}
	fmt.Printf("%v", s.Router)
	http.ListenAndServe(":10101", nil)
	return nil
}
