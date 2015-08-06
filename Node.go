package main

import (
	"encoding/json"
)

type TsNode struct {
	input   chan TsPacket
	outputs []chan<- TsPacket
	PktsIn  int64 //counters
	PktsOut int64
}

// accessor to get this node's input channel
func (node *TsNode) GetInputChan() chan TsPacket {
	return node.input
}

// add a new listener to the output list
func (node *TsNode) RegisterListener(newout chan<- TsPacket) {
	node.outputs = append(node.outputs, newout)
}

// remove a particular listener from the outputs slice
// does nothing if the chan is not registered
func (node *TsNode) UnRegisterListener(toremove chan<- TsPacket) {
	for i, op := range node.outputs {
		if op == toremove {
			node.outputs = append(node.outputs[:i], node.outputs[i+1:]...)
			break
		}
	}
}

// dump a representation of this node to JSON
// this could be used by a web interface etc to monitor the status
// of each node. This method may be hidden by a struct which embeds
// TsNode, in which case that struct will also dump its own data too.
func (node *TsNode) ToJson() ([]byte, error) {
	return json.Marshal(node)
}
