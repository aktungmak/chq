package main

import (
	"errors"
)

type Routeable interface {
	RegisterListener(chan<- TsPacket)
	UnRegisterListener(chan<- TsPacket)
	GetInputChan() chan TsPacket
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

func (r *Router) RegisterNode(name string, newnode Routeable) {
	r.Nodes[name] = newnode
}

func (r *Router) Connect(src, dst string) error {
	sn, ok := r.Nodes[src]
	dn, ok := r.Nodes[dst]
	if !ok {
		return errors.New("No such node !")
	}
	sn.RegisterListener(dn.GetInputChan())
	return nil
}

func (r *Router) Disconnect(src, dst string) error {
	sn, ok := r.Nodes[src]
	dn, ok := r.Nodes[dst]
	if !ok {
		return errors.New("No such node !")
	}
	sn.UnRegisterListener(dn.GetInputChan())
	return nil
}

//should there be a disconnect function?
//there would need to be UnRegisterNode too
