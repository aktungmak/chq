package main

import (
	"log"
)

// A PidFilter allows only the specified PID to pass
// all other PIDs will be rejected.
// compare with PidDropper, which has the opposite behaviour
type PidFilter struct {
	Pid int16
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("PidFilter", NewPidFilter)

}

func NewPidFilter(pid int16) (*PidFilter, error) {
	node := &PidFilter{}
	node.Pid = pid
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)
	node.outputs = make([]chan<- TsPacket, 0)

	go node.process()
	return node, nil
}

func (node *PidFilter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		if pkt.Header.Pid == node.Pid {
			for _, output := range node.outputs {
				output <- pkt
			}
		}
	}
}

func (node *PidFilter) closeDown() {
	log.Printf("closing down pid filter for pid %d", node.Pid)
	for _, output := range node.outputs {
		close(output)
	}
}
