package main

import (
	"fmt"
	"log"
)

// A CcCommenter will watch the CC field of
// all TSPackets it receives. If it notices
// a CC error, it will log this in the
// Comment field of the TSPacket struct.
// It passes through all packets received.
type CcCommenter struct {
	curCc map[int]byte
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("CcCommenter", NewCcCommenter)
}

func NewCcCommenter() (*CcCommenter, error) {
	node := &CcCommenter{}
	node.curCc = make(map[int]byte)
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *CcCommenter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		node.PktsIn++
		node.Send(pkt)
		// filter out NULL pid and erroneous values
		if pkt.Header.Pid < 0x1FF {
			prev, ok := node.curCc[pkt.Header.Pid]
			node.curCc[pkt.Header.Pid] = pkt.Header.Cc

			// have we seen this one before?
			if ok {
				if pkt.Header.Cc != ((prev + 1) % 16) {
					// no need to inc if no payld
					if pkt.Header.Afc&1 == 1 {
						pkt.Comment = fmt.Sprintf("CC error: expected %d got %d (packet %d)",
							prev+1,
							pkt.Header.Cc,
							node.PktsIn)
					}
				}
			}
		}
	}
}

func (node *CcCommenter) closeDown() {
	log.Print("closing down CcCommenter")
	node.output.Close()
}
