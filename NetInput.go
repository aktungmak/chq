package main

import (
	"encoding/json"
	"log"
	"net"
)

type NetInput struct {
	addr net.UDPAddr
	conn *net.UDPConn
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes["NetInput"] = NewNetInput
}

func NewNetInput(address string, port int) (*NetInput, error) {
	var err error
	var conn *net.UDPConn

	n := &NetInput{}
	n.addr = net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(address),
	}
	n.input = nil
	n.outputs = make([]chan<- TsPacket, 0)

	// the stream may be unicast or multicast, so choose appropriately
	if n.addr.IP.IsMulticast() {
		conn, err = net.ListenMulticastUDP("udp", nil, &n.addr)
	} else {
		conn, err = net.ListenUDP("udp", &n.addr)
	}
	if err != nil {
		conn.Close()
		return nil, err
	}
	n.conn = conn

	go n.process()
	return n, nil
}

func (node *NetInput) ToJson() ([]byte, error) {
	return json.Marshal(node)
}

func (node *NetInput) process() {
	defer node.closeDown()
	var packetsize int
	packet := make([]byte, 4096)
	for {
		m := 0
		n, _, err := node.conn.ReadFromUDP(packet)
		if err != nil {
			log.Printf("TS capture error: %s", err)
			continue
		}
		log.Printf("got %d bytes", n)
		if (packet[0] & 192) == 128 {
			// this is RTP, skip the header
			m += 12 + (4 * (int(packet[0]) & 15))
		}

		// check packetsize (188/204)
		if ((n - m) % 188) == 0 {
			packetsize = 188
		} else if ((n - m) % 204) == 0 {
			// 204 bytes pkts
			packetsize = 204
		} else {
			panic("Unknown TS packet size!!")
		}
		log.Printf("packetsize is %d bytes", packetsize)
		// split into TS packets
		for i := m; i < n; i += packetsize {
			pkt := NewTsPacket(packet[i : i+packetsize])
			node.PktsIn++
			for _, output := range node.outputs {
				node.PktsOut++
				output <- pkt
			}
		}
	}
}

func (node *NetInput) closeDown() {
	node.conn.Close()
	for _, output := range node.outputs {
		close(output)
	}
}
