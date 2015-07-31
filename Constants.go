package main

const (
	CHAN_BUF_SIZE   = 10
	TS_PKT_SIZE     = 188
	FILE_CHUNK_SIZE = 500 //number of packets to read at a time
)

// all nodes should register here so they can be looked up
// TODO make this a type and give it a Register() method to
// prevent Node types from hiding each other.
var AvailableNodes = make(map[string]interface{})
