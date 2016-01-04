package main

import (
	"log"
)

// PcrBrCommenter measures the bitrate of the
// stream passing through it using the specified
// PCR PID. This is an "offline" method. The current
// bitrate will be recorded in the packet comment.
type PcrBrCommenter struct {
	PcrPid  int
	LastPcr int64
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("PcrBrCommenter", NewPcrBrCommenter)
}

func NewPcrBrCommenter(PcrPid, int) (*PcrBrCommenter, error) {
	node := &PcrBrCommenter{}
	node.PcrPid = PcrPid
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *PcrBrCommenter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		node.PktsIn++
		if pkt.Header.Afc && pkt.AdaptationField.Pcrf {

		}
		node.Send(pkt)
	}
}

func (node *PcrBrCommenter) closeDown() {
	log.Printf("closing down PcrBrCommenter on pcr pid %d", node.PcrPid)
	node.output.Close()
}
