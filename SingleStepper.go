package main

import (
	"log"
)

// SingleStepper only lets one packet through
// before becoming inactive again.
// It can be used to step through a TS slowly
type SingleStepper struct {
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("SingleStepper", NewSingleStepper)
}

func NewSingleStepper() (*SingleStepper, error) {
	node := &SingleStepper{}
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *SingleStepper) process() {
	defer node.closeDown()
	for pkt := range node.input {
		node.PktsIn++
		node.Send(pkt)

		node.L.Lock()
		node.Active = false
		node.L.Unlock()
	}
}

func (node *SingleStepper) closeDown() {
	node.Active = false
	log.Print("closing down SingleStepper")
	node.output.Close()
}
