package main

import (
	"log"
)

// PidRemapper changes the PID of all incoming
// packets to ToPid if they match FromPid
// All other packets are passed through
type PidRemapper struct {
	FromPid int
	ToPid   int
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("PidRemapper", NewPidRemapper)
}

func NewPidRemapper(FromPid, ToPid int) (*PidRemapper, error) {
	node := &PidRemapper{}
	node.FromPid = FromPid
	node.ToPid = ToPid
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *PidRemapper) process() {
	defer node.closeDown()
	for pkt := range node.input {
		node.PktsIn++
		if pkt.Header.Pid != node.FromPid {
			node.Send(pkt)
		} else {
			// copy the bytes
			dat := append([]byte{}, pkt.bytes...)

			// edit the pid data
			dat[1] = byte((node.ToPid>>8)&0x31) | (dat[1] & 0xe0)
			dat[2] = byte(node.ToPid)

			opkt := NewTsPacket(dat)
			node.Send(opkt)
		}
	}
}

func (node *PidRemapper) closeDown() {
	node.Active = false
	log.Printf("closing down pidremapper from pid %d to pid %d", node.FromPid, node.ToPid)
	node.output.Close()
}
