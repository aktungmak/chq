package main

import (
	"encoding/json"
	"log"
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
	control struct {
		active bool
		c      sync.Cond
	}
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

// send the provided packet using our output
// broadcaster. if we are not active, wait for
// the signal. increments counters appropriately.
func (node *TsNode) Send(pkt TsPacket) {
	for !node.control.active {
		node.control.c.Wait()
	}
	node.output.Send(pkt)
	node.PktsOut++
}

// Switch a node between being active/inactive states
// not all nodes used this, but it is a good idea for
// sources to use this so they don't start outputting
// before downstream is ready.
func (node *TsNode) Toggle() {
	log.Print("Togggggle!")
	node.control.active = !node.control.active
	node.control.c.Signal()
}

func (node *TsNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		PktsIn  int64
		PktsOut int64
		Active  bool
	}{
		node.PktsIn,
		node.PktsOut,
		node.control.active,
	})
}
