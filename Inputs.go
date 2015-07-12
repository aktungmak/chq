package main

import (
	"log"
	"net"
	"os"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}

}

func FileInput(fname string) (<-chan TsPacket, error) {
	var err error
	output := make(chan TsPacket, CHAN_BUF_SIZE)

	// try to open file
	f, err := os.Open(fname)
	Check(err)

	buf := make([]byte, TS_PKT_SIZE)

	go func() {
		defer f.Close()
		for {
			// check for sync
			n, err := f.Read(buf)
			Check(err)
			if n < TS_PKT_SIZE {
				log.Printf("Couldn't get a full packet, only %d bytes", n)
				continue
			} else if buf[0] != 0x47 {
				log.Print("no lock yet")
				f.Seek(-int64(n-1), 1)
				continue
			} else {
				output <- NewTsPacket(buf)
			}

		}
	}()
	return output, nil
}

func NetInput(address string, port int) (<-chan TsPacket, error) {
	var err error
	var conn *net.UDPConn
	var packetsize int
	output := make(chan TsPacket, CHAN_BUF_SIZE)

	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(address),
	}

	// the stream may be unicast or multicast, so choose appropriately
	if addr.IP.IsMulticast() {
		conn, err = net.ListenMulticastUDP("udp", nil, &addr)
	} else {
		conn, err = net.ListenUDP("udp", &addr)
	}

	if err != nil {
		conn.Close()
		return nil, err
	}

	defer conn.Close()
	packet := make([]byte, 4096)
	for {
		m := 0
		n, _, err := conn.ReadFromUDP(packet)
		if err != nil {
			log.Printf("TS capture error: %s", err)
			continue
		}
		log.Printf("got %d packets", n)
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
			output <- NewTsPacket(packet[i : i+packetsize])
		}
	}
	return output, nil
}
