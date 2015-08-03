package main

import (
	"errors"
)

type Routeable interface {
	RegisterListener(chan<- TsPacket)
	UnRegisterListener(chan<- TsPacket)
	GetInputChan() chan TsPacket
	// ToJson() []byte
}

type Router struct {
	Nodes map[string]Routeable
}

// construct a new, blank router
func NewRouter() *Router {
	return &Router{
		Nodes: make(map[string]Routeable),
	}
}

func (r *Router) RegisterNode(name string, newnode Routeable) error {
	_, present := r.Nodes[name]
	if present {
		return errors.New("Node already exists: " + name)
	} else {
		r.Nodes[name] = newnode
		return nil
	}
}

func (r *Router) Connect(src, dst string) error {
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
