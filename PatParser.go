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
	AvailableNodes.Register("PatParser", NewPatParser)
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
		node.PktsIn++
		if pkt.Header.Pid == 0 {

			//if pusi
			if pkt.Header.Pusi {
				ptr := int(pkt.Payload[0])
				if secLen == 0 {
					// push data from the ptr onwards into the buffer
					copy(secBuf[secLen:], pkt.Payload[ptr+1:])
					secLen += len(pkt.Payload) - ptr - 1
				} else { //(buf > 0)

					//        push data up until the ptr
					//        increment bufLen
					//        parse buffer
					//
					//        set bufLen to 0
					//        push data after the ptr
					//
				}
			} else { //no pusi
				if secLen > 0 {
					if secLen+len(pkt.Payload) > 4096 {
						log.Print("PAT has overflowed the section buffer!")
						secLen = 0
						continue
					} else {
						//        push data into buffer
						//        increment bufLen
					}
				}

			}

			for _, output := range node.outputs {
				node.PktsOut++
				output <- pkt
			}
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
