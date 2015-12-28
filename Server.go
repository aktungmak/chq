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

// return the info on a particular node, specified by name
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

// list out all the available node types
// TODO make into a map of name: description
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

// on a GET, list out all the connections for that node
// on a POST, make a new connection
// TODO add DELETE to disconnect
func (s *Server) conn(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/conn/"):]
	node, ok := s.Router.Nodes[name]
	if !ok {
		http.Error(w, "node not found!", http.StatusNotFound)
	}
	switch r.Method {

	case "GET":
		res := make([]string, 0)
		outs := node.GetOutputs()
		for n, a := range s.Router.Nodes {
			for _, b := range outs {
				if a.GetInputChan() == b {
					res = append(res, n)
				}
			}
		}
		dat, err := marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write(dat)
			return
		}

	case "POST":
		// TODO implement making a new connection
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) state(w http.ResponseWriter, r *http.Request) {
	// todo implement state toggling for each node or all nodes
}

// start a server serving. blocks until err or exit
func (s *Server) Start() error {
	var handlers = map[string]func(http.ResponseWriter, *http.Request){
		"/":        s.summary,
		"/summary": s.summary,
		"/nodes/":  s.nodes,
		"/types":   s.types,
		"/conn/":   s.conn,
		"/state/":  s.state,
	}

	for ep, h := range handlers {
		http.HandleFunc(ep, h)
	}
	fmt.Printf("%v", s.Router)
	http.ListenAndServe(":10101", nil)
	return nil
}
