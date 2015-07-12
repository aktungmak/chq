package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	log.Print("starting")
	// fname := flag.String("if", "", "input file")
	addr := flag.String("a", "", "input multicast/unicast address")
	port := flag.Int("p", 0, "input UDP port")
	flag.Parse()
	// op, err := FileInput(*fname)
	op, err := NetInput(*addr, *port)
	Check(err)
	fmt.Printf("%v", <-op)
}
