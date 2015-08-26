package main

import (
	"log"
)

const (
	CHAN_BUF_SIZE   = 10
	TS_PKT_SIZE     = 188
	FILE_CHUNK_SIZE = 500 //number of packets to read at a time
	PAT_PID         = 0
)

type AvailableNodeMap map[string]interface{}

// TODO convert all nodes to use this method
func (a AvailableNodeMap) Register(name string, node interface{}) {
	_, exist := a[name]
	if exist {
		log.Fatalf("Node '%s' already registered!", name)
	} else {
		a[name] = node
	}
}

// all nodes should register here so they can be looked up
var AvailableNodes = make(AvailableNodeMap)
