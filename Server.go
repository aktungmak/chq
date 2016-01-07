package main

import (
	"encoding/json"
	"net/http"
	"strings"
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
	v, err := s.Router.GetNodeByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		dat, err := marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			w.Write(dat)
		}
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
	segs := strings.Split(r.URL.Path, "/")
	node, err := s.Router.GetNodeByName(segs[2])
	if err != nil {
		http.Error(w, "node not found!", http.StatusNotFound)
	} else {

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
				return
			} else {
				w.Write(dat)
				return
			}

		case "POST":
			// expect uri like /conn/src/dst
			if len(segs) < 3 {
				http.Error(w, "src & dest nodes not specified!", http.StatusNotFound)
				return
			}
			src := segs[2]
			dst := segs[3]
			err = s.Router.Connect(src, dst)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusCreated)
				return
			}

		case "DELETE":
			// expect uri like /conn/src/dst
			if len(segs) < 3 {
				http.Error(w, "src & dest nodes not specified!", http.StatusNotFound)
				return
			}
			src := segs[2]
			dst := segs[3]
			err := s.Router.Disconnect(src, dst)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			} else {
				w.WriteHeader(http.StatusNoContent)
				return
			}

		default:
			http.Error(w, "", http.StatusMethodNotAllowed)
		}
	}
}

func (s *Server) state(w http.ResponseWriter, r *http.Request) {
	// todo implement state toggling for each node or all nodes
	// PUT to /state/ toggles all nodes
	// PUT request to /state/nodename will trigger a toggle
	// response with bool of node's new state
	name := r.URL.Path[len("/state/"):]
	if len(name) == 0 {
		s.Router.ToggleAll()
		return
	}
	node, err := s.Router.GetNodeByName(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else {
		// todo return new state
		node.Toggle()
	}

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
	http.ListenAndServe(":10101", nil)
	return nil
}
