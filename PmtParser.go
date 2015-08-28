package main

import (
	"encoding/json"
	"log"
)

// PmtParser watches the incoming TS packets for
// its Pid. It will try to reassemble these into a
// complete section, and then parse as a Pmt.
// When the version number changes, it will push
// the CurPmt onto the slice of PrevPmts.
type PmtParser struct {
	PrevPmts []*Pmt
	CurPmt   *Pmt
	Pid      int16
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("PmtParser", NewPmtParser)
}

func NewPmtParser(pid int16) (*PmtParser, error) {
	node := &PmtParser{}
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)
	node.outputs = make([]chan<- TsPacket, 0)

	node.PrevPmts = make([]*Pmt, 0)
	node.Pid = pid

	go node.process()
	return node, nil
}

func (node *PmtParser) process() {
	defer node.closeDown()
	// section_length is 12 bits value
	secBuf := make([]byte, 4096)
	bufLen := 0
	for pkt := range node.input {
		node.PktsIn++
		if pkt.Header.Pid == node.Pid {
			if pkt.Header.Pusi {
				ptr := int(pkt.Payload[0])
				if bufLen > 0 {
					// push data up to the ptr into the buffer
					copy(secBuf[bufLen:], pkt.Payload[1:ptr+1])
					bufLen += len(pkt.Payload) - ptr - 1 // IS IT -1???
					pmt, err := NewPmt(secBuf[:bufLen])
					if err != nil {
						log.Print(err)
					} else {
						if node.CurPmt != nil && node.CurPmt.Vn != pmt.Vn {
							node.PrevPmts = append(node.PrevPmts, node.CurPmt)
						}
						node.CurPmt = pmt
					}
					bufLen = 0

				}

				// copy the next section into the buffer
				copy(secBuf[bufLen:], pkt.Payload[ptr+1:])
				bufLen += len(pkt.Payload) - ptr - 1

			} else {
				if bufLen > 0 {
					if bufLen+len(pkt.Payload) > 4096 {
						log.Print("PMT has overflowed the section buffer!")
						bufLen = 0
						continue
					} else {
						copy(secBuf[bufLen:], pkt.Payload)
						bufLen += len(pkt.Payload)
					}
				}

			}

		}
		for _, output := range node.outputs {
			node.PktsOut++
			output <- pkt
		}
	}
}

func (node *PmtParser) closeDown() {
	log.Print("closing down PmtParser")
	for _, output := range node.outputs {
		close(output)
	}
}

func (node *PmtParser) ToJson() ([]byte, error) {
	return json.Marshal(node)
}
