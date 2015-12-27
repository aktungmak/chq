package main

import (
	"encoding/json"
	"log"
	"strconv"
)

// PidCounter keeps a count of all the pids
// that it has seen so far.
type PidCounter struct {
	Pids map[int]int
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("PidCounter", NewPidCounter)
}

func NewPidCounter() (*PidCounter, error) {
	node := &PidCounter{}
	node.Pids = make(map[int]int)
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *PidCounter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		node.PktsIn++
		node.Send(pkt)

		node.Pids[pkt.Header.Pid]++
	}
}

func (node *PidCounter) closeDown() {
	log.Printf("closing down pid counter, found %d pids", len(node.Pids))
	node.output.Close()
}

func (node *PidCounter) MarshalJSON() ([]byte, error) {
	// make a dumme struct that looks like the node
	// but is suitable for json
	tmp := struct {
		TsNode
		Pids map[string]int
	}{
		node.TsNode,
		make(map[string]int),
	}

	for pid, cnt := range node.Pids {
		tmp.Pids[strconv.Itoa(pid)] = cnt
	}

	return json.Marshal(tmp)
}
