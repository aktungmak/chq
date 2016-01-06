package main

import (
	"fmt"
	"log"
)

// PcrBrCommenter measures the bitrate of the
// stream passing through it using the specified
// PCR PID. This is an "offline" method. The current
// bitrate will be recorded in the packet comment.
type PcrBrCommenter struct {
	PcrPid     int
	LastPcr    int64
	lastPktCnt int64
	CurBr      float64
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("PcrBrCommenter", NewPcrBrCommenter)
}

func NewPcrBrCommenter(pcrPid int) (*PcrBrCommenter, error) {
	node := &PcrBrCommenter{}
	node.PcrPid = pcrPid
	node.LastPcr = -MAX_PCR_STEP - 1 // sentinel to detect first run
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *PcrBrCommenter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		node.PktsIn++
		if pkt.Header.Pid == node.PcrPid && (pkt.Header.Afc > 1) && pkt.AdaptationField.Pcrf {
			// todo: don't just use the PCR base, use ext also
			dPkt := float64(node.PktsIn - node.lastPktCnt)
			dPcr := float64(pkt.AdaptationField.Pcrb - node.LastPcr)

			if dPcr > MAX_PCR_STEP {
				// if last PCR is negative, this is first run so discon expected
				if node.LastPcr > 0 {
					// pcr discon, resync
					log.Printf("PCR jumps by more than 10 sec (%.0f ticks) in packet %d", dPcr, node.PktsIn)
				}
			} else {
				node.CurBr = (dPkt * TS_PKT_SIZE * 8) / (dPcr / 90000.0)
				pkt.Comment = fmt.Sprintf("%f", node.CurBr)
			}

			node.LastPcr = pkt.AdaptationField.Pcrb
			node.lastPktCnt = node.PktsIn
		}
		node.Send(pkt)
	}
}

func (node *PcrBrCommenter) closeDown() {
	log.Printf("closing down PcrBrCommenter on pcr pid %d", node.PcrPid)
	node.output.Close()
}
