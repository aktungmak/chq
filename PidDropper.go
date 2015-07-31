package main

import (
	"log"
)

// A PidDropper allows all but the specified PID to pass
// the selected PID will be dropped
// compare with PidFilter, which has the opposite behaviour
type PidDropper struct {
	Pid Pid
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes["PidDropper"] = NewPidDropper
}

func NewPidDropper(pid Pid) (*PidDropper, error) {
	node := &PidDropper{}
	node.Pid = pid
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)
	node.outputs = make([]chan<- TsPacket, 0)

	go node.process()
	return node, nil
}

func (node *PidDropper) process() {
	defer node.closeDown()
	for pkt := range node.input {
		if pkt.Header.Pid != node.Pid {
			for _, output := range node.outputs {
				output <- pkt
			}
		}
	}
}

func (node *PidDropper) closeDown() {
	log.Printf("closing down pid dropper for pid %d", node.Pid)
	for _, output := range node.outputs {
		close(output)
	}
}
