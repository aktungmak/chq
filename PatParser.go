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

			if pkt.Header.Pusi {
				ptr := int(pkt.Payload[0])
				if ptr > 0 {
					log.Printf("%v", pkt.Payload)
				}
				// if overflow, give up
				if secLen+ptr > len(secBuf) {
					log.Print("PAT has overflowed the section buffer!")
					secLen = 0
					// otherwise, if we already have data, get the last and parse
				} else if secLen > 0 {
					// copy first half to buffer
					copy(secBuf[secLen:], pkt.Payload[1:ptr+1])
					// inc secLen
					secLen += ptr
					// parse buffer
					np, err := NewPat(secBuf[:secLen])
					if err != nil {
						log.Printf("Error parsing PAT: %s", err)
					} else {
						node.CurPat = np
					}
				}
				// clear buffer
				secLen = 0

				// starting a new PAT here with the rest of the data
				copy(secBuf[secLen:], pkt.Payload[ptr:])
				// copy(secBuf[secLen:], pkt.Payload[ptr+1:])

			} else {
				// this is just PAT payload
				// check for overflow
				if secLen+len(pkt.Payload) > len(secBuf) {
					log.Print("PAT has overflowed the section buffer!")
					secLen = 0
				} else {
					// all is ok, so push it into the buffer
					copy(secBuf[secLen:], pkt.Payload)
					secLen += len(pkt.Payload)
				}
			}

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

func (node *PatParser) ToJson() ([]byte, error) {
	return json.Marshal(node)
}
