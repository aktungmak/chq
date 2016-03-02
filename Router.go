package main

import (
	"errors"
	"sync"
)

type Routeable interface {
	GetInputChan() chan TsPacket
	RegisterListener(chan TsPacket)
	UnRegisterListener(chan TsPacket)
	GetOutputs() []chan TsPacket
	Toggle() bool
}

type Router struct {
	Nodes map[string]Routeable
	sync.Mutex
}

// construct a new, blank router
func NewRouter() *Router {
	return &Router{
		Nodes: make(map[string]Routeable),
	}
}

func (r *Router) GetNodeByName(name string) (Routeable, error) {
	r.Lock()
	defer r.Unlock()
	node, ok := r.Nodes[name]
	if !ok {
		return nil, errors.New("No such node " + name)
	} else {
		return node, nil
	}
}

func (r *Router) RegisterNode(name string, newnode Routeable) error {
	r.Lock()
	defer r.Unlock()
	_, present := r.Nodes[name]
	if present {
		return errors.New("Node already exists: " + name)
	} else {
		r.Nodes[name] = newnode
		return nil
	}
}

func (r *Router) Connect(src, dst string) error {
	r.Lock()
	defer r.Unlock()
	sn, ok := r.Nodes[src]
	if !ok {
		return errors.New("No such node " + src)
	}
	dn, ok := r.Nodes[dst]
	if !ok {
		return errors.New("No such node " + dst)
	}
	sn.RegisterListener(dn.GetInputChan())
	return nil
}

func (r *Router) Disconnect(src, dst string) error {
	r.Lock()
	defer r.Unlock()
	sn, ok := r.Nodes[src]
	if !ok {
		return errors.New("No such node " + src)
	}
	dn, ok := r.Nodes[dst]
	if !ok {
		return errors.New("No such node " + dst)
	}
	sn.UnRegisterListener(dn.GetInputChan())
	return nil
}

// trigger all nodes in the router. should be
// used after ApplyConfig is complete to start
// all nodes processing.
func (r *Router) ToggleAll() {
	r.Lock()
	defer r.Unlock()
	for _, node := range r.Nodes {
		node.Toggle()
	}
}
