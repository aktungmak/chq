package main

import (
	"sync"
)

// TsNode is the core of every node. It provides the
// inputs and outputs used by every type of node, and
// defines the methods specified by Routeable. It is
// must be embedded in any new Node struct.
type TsNode struct {
	input   chan TsPacket
	output  Broadcaster
	PktsIn  int64
	PktsOut int64
	// Active controls whether we are processing or not
	Active bool
	// this Cond syncs access to Active
	// its .L field is lazily initialized
	sync.Cond `json:"-"`
}

// accessor to get this node's input channel
func (node *TsNode) GetInputChan() chan TsPacket {
	return node.input
}

// add a new listener to the output list
func (node *TsNode) RegisterListener(newout chan TsPacket) {
	node.output.RegisterChan(newout)
}

// remove a particular listener from the outputs slice
// does nothing if the chan is not registered
func (node *TsNode) UnRegisterListener(toremove chan TsPacket) {
	node.output.UnRegisterChan(toremove)
}

// list out all the currently available outputs
func (node *TsNode) GetOutputs() []chan TsPacket {
	return node.output.GetOutputs()
}

// send the provided packet using our output
// broadcaster. if we are not active, wait for
// the signal. increments counters appropriately.
func (node *TsNode) Send(pkt TsPacket) {
	if node.L == nil {
		node.Cond = *sync.NewCond(&sync.Mutex{})
	}
	for !node.Active {
		node.L.Lock()
		node.Wait()
	}
	node.output.Send(pkt)
	node.PktsOut++
}

// Switch a node between being active/inactive states
// The Send() method will block until node.Active == true
func (node *TsNode) Toggle() {
	if node.L == nil {
		node.Cond = *sync.NewCond(&sync.Mutex{})
	}
	node.L.Lock()
	node.Active = !node.Active
	node.L.Unlock()
	node.Signal()
}
