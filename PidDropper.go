package main

import (
	"log"
)

// A PidDropper allows all but the specified PID to pass
// the selected PID will be dropped
// compare with PidFilter, which has the opposite behaviour
type PidDropper struct {
	Pid int
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("PidDropper", NewPidDropper)
}

func NewPidDropper(pid int) (*PidDropper, error) {
	node := &PidDropper{}
	node.Pid = pid
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *PidDropper) process() {
	defer node.closeDown()
	for pkt := range node.input {
		if pkt.Header.Pid != node.Pid {
			node.PktsOut++
			node.output.Send(pkt)
		}
	}
}

func (node *PidDropper) closeDown() {
	log.Printf("closing down pid dropper for pid %d", node.Pid)
	node.output.Close()
}
