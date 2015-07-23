package main

type TsNode struct {
	input   chan TsPacket
	outputs []chan<- TsPacket
	pktsIn  int //counters
	pktsOut int
}

// return this node's input channel
func (node *TsNode) GetInputChan() chan TsPacket {
	return node.input
}

// add a new listener to the output list
func (node *TsNode) RegisterListener(newout chan<- TsPacket) {
	node.outputs = append(node.outputs, newout)
}

func (node *TsNode) UnRegisterListener(toremove chan<- TsPacket) {
	for i, op := range node.outputs {
		if op == toremove {
			node.outputs = append(node.outputs[:i], node.outputs[i+1:]...)
			break
		}
	}
}
