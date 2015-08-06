package main

import (
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
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes["PatParser"] = NewPatParser
}

func NewPatParser(pid int16) (*PatParser, error) {
	node := &PatParser{}
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)
	node.outputs = make([]chan<- TsPacket, 0)

	go node.process()
	return node, nil
}

func (node *PatParser) process() {
	defer node.closeDown()
	// section_length is 12 bits value
	secBuf := make([]byte, 4096)
	secLen := 0
	for pkt := range node.input {
		if pkt.Header.Pid == 0 {
			if (secLen + len(pkt.Payload)) > len(secBuf) {
				log.Print("PAT has overflowed the section buffer!")
				secLen = 0
			} else {
				copy(secBuf[secLen], pkt.Payload)
			}

			if pkt.Header.Pusi {
				if secLen > 0 {
					// this completes the previous section
					node.CurPat = NewPat(data[:secLen])
				}
			}
			//
		}
		for _, output := range node.outputs {
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
