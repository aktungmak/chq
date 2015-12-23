package main

import (
	"log"
)

const (
	CHAN_BUF_SIZE   = 1
	TS_PKT_SIZE     = 188
	FILE_CHUNK_SIZE = 500 //number of packets to read at a time
)

type AvailableNodeMap map[string]interface{}

// all nodes must use use this method to make themselves available
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
