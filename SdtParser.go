package main

import (
	"log"
)

// SdtParser watches the incoming TS packets for
// its Pid. It will try to reassemble these into a
// complete section, and then parse as an Sdt.
// When the version number changes, it will push
// the CurSdt onto the slice of PrevSdts.
type SdtParser struct {
	PrevSdts []*Sdt
	CurSdt   *Sdt
	Pid      int
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("SdtParser", NewSdtParser)
}

func NewSdtParser(pid int) (*SdtParser, error) {
	node := &SdtParser{}
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	node.PrevSdts = make([]*Sdt, 0)
	node.Pid = pid

	go node.process()
	return node, nil
}

func (node *SdtParser) process() {
	defer node.closeDown()
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
					Sdt, err := NewSdt(secBuf[:bufLen])
					if err != nil {
						log.Print(err)
					} else {
						if node.CurSdt != nil && node.CurSdt.Vn != Sdt.Vn {
							node.PrevSdts = append(node.PrevSdts, node.CurSdt)
						}
						node.CurSdt = Sdt
					}
					bufLen = 0

				}

				// copy the next section into the buffer
				copy(secBuf[bufLen:], pkt.Payload[ptr+1:])
				bufLen += len(pkt.Payload) - ptr - 1

			} else {
				if bufLen > 0 {
					if bufLen+len(pkt.Payload) > 4096 {
						log.Print("SDT has overflowed the section buffer!")
						bufLen = 0
						continue
					} else {
						copy(secBuf[bufLen:], pkt.Payload)
						bufLen += len(pkt.Payload)
					}
				}

			}

		}
		node.Send(pkt)
	}
}

func (node *SdtParser) closeDown() {
	node.Active = false
	log.Printf("closing down SdtParser for pid %d", node.Pid)
	node.output.Close()
}
