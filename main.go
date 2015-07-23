package main

import (
	"flag"
	"log"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	log.Print("starting")
	fname := flag.String("if", "", "input file")

	// addr := flag.String("a", "", "input multicast/unicast address")
	// port := flag.Int("p", 0, "input UDP port")
	flag.Parse()

	rtr := NewRouter()

	fin, err := NewFileInput(*fname)
	pdr, err := NewPidDropper(17)
	pdq, err := NewPidDropper(0)
	// op, err := NewNetInput(*addr, *port)

	rtr.RegisterNode("filein", fin)
	rtr.RegisterNode("pidrup", pdq)
	rtr.RegisterNode("pidrop", pdr)

	Check(err)
	inp := make(chan TsPacket)
	rtr.Connect("filein", "pidrup")
	rtr.Connect("pidrup", "pidrop")

	// rtr.Disconnect("filein", "pidrup")
	pdr.RegisterListener(inp)

	pids := make(map[Pid]int)
	for pkt := range inp {
		pids[pkt.Header.Pid]++
	}
	log.Printf("%v", pids)
	// pkt = <-inp
	// log.Printf("got: %v", pkt.Header.Pid)
	// pkt = <-inp
	// log.Printf("got: %v", pkt.Header.Pid)
	// pkt = <-inp
	// log.Printf("got: %v", pkt.Header.Pid)
}
