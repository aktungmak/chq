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
	CurCc map[int16]byte
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("CcCommenter", NewCcCommenter)
}

func NewCcCommenter() (*CcCommenter, error) {
	node := &CcCommenter{}
	node.CurCc = make(map[int16]byte)
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *CcCommenter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		// filter out NULL pid and erroneous values
		if pkt.Header.Pid < 0x1FF {
			prev, ok := node.CurCc[pkt.Header.Pid]
			node.CurCc[pkt.Header.Pid] = pkt.Header.Cc

			// have we seen this one before?
			if ok {
				if pkt.Header.Cc != ((prev + 1) % 16) {
					// no need to inc if no payld
					if pkt.Header.Afc&1 == 1 {
						pkt.Comment = fmt.Sprintf("CC error: expected %d got %d",
							prev+1,
							pkt.Header.Cc)
					}
				}
			}
		}

		node.output.Send(pkt)
	}
}

func (node *CcCommenter) closeDown() {
	log.Print("closing down CcCommenter")
	node.output.Close()
}
