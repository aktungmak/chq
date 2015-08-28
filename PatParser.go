package main

import (
	"encoding/json"
	"log"
)

// PatParser watches the incoming TS packets for
// PID 0. It will try to reassemble these into a
// complete section, and then parse as a PAT.
// When the version number changes, it will push
// the CurPat onto the slice of PrevPats.
type PatParser struct {
	PrevPats []*Pat
	CurPat   *Pat
	Pid      int16
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("PatParser", NewPatParser)
}

func NewPatParser(pid int16) (*PatParser, error) {
	node := &PatParser{}
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)
	node.outputs = make([]chan<- TsPacket, 0)

	node.PrevPats = make([]*Pat, 0)
	node.Pid = pid

	go node.process()
	return node, nil
}

func (node *PatParser) process() {
	defer node.closeDown()
	// section_length is 12 bits value
	secBuf := make([]byte, 4096)
	bufLen := 0
	for pkt := range node.input {
		node.PktsIn++
		if pkt.Header.Pid == node.Pid {
			if pkt.Header.Pusi { //yes pusi DONE
				ptr := int(pkt.Payload[0])
				if bufLen > 0 {
					// push data up to the ptr into the buffer
					copy(secBuf[bufLen:], pkt.Payload[1:ptr+1])
					bufLen += len(pkt.Payload) - ptr - 1 // IS IT -1???
					pat, err := NewPat(secBuf[:bufLen])
					if err != nil {
						log.Print(err)
					} else {
						if node.CurPat != nil && node.CurPat.Vn != pat.Vn {
							node.PrevPats = append(node.PrevPats, node.CurPat)
						}
						node.CurPat = pat
					}
					bufLen = 0

				}

				// copy the next section into the buffer
				copy(secBuf[bufLen:], pkt.Payload[ptr+1:])
				bufLen += len(pkt.Payload) - ptr - 1

			} else { //no pusi DONE
				if bufLen > 0 {
					if bufLen+len(pkt.Payload) > 4096 {
						log.Print("PAT has overflowed the section buffer!")
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

func (node *PatParser) closeDown() {
	log.Print("closing down PatParser")
	for _, output := range node.outputs {
		close(output)
	}
}

func (node *PatParser) ToJson() ([]byte, error) {
	return json.Marshal(node)
}
