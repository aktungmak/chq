package main

import (
	"encoding/json"
)

// TsNode is the core of every node. It provides the
// inputs and outputs used by every type of node, and
// defines the methods specified by Routeable. It is
// must be embedded in any new Node struct.
type TsNode struct {
	input   chan TsPacket
	output  Broadcaster
	PktsIn  int64 //counters
	PktsOut int64
}

// // read a packet off the input. second value is
// // false if the input channel is closed.
// func (node *TsNode) GetPacket() (TsPacket, bool) {
// 	data, ok := <-node.input
// 	pkt := TsPacket(data)
// 	return pkt, ok
// }

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

// dump a representation of this node to JSON
// this could be used by a web interface etc to monitor the status
// of each node. This method may be hidden by a struct which embeds
// TsNode, in which case that struct will also dump its own data too.
func (node *TsNode) ToJson() ([]byte, error) {
	return json.Marshal(node)
}
