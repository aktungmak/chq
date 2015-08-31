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

	go node.process()
	return node, nil
}

func (node *PidFilter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		if pkt.Header.Pid == node.Pid {
			node.output.Send(pkt)
		}
	}
}

func (node *PidFilter) closeDown() {
	log.Printf("closing down PidFilter for pid %d", node.Pid)
	node.output.Close()
}
