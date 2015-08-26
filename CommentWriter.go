package main

import (
	"log"
	"os"
)

// A CommentWriter will append the comment fields
// of all TS packets it receives to a file.
// It passes through TS packets unmodified.
type CommentWriter struct {
	file *os.File
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("CommentWriter", NewCommentWriter)
}

func NewCommentWriter(fname string) (*CommentWriter, error) {
	// try to open file
	fh, err := os.OpenFile(fname, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	node := &CommentWriter{}
	node.file = fh
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)
	node.outputs = make([]chan<- TsPacket, 0)

	go node.process()
	return node, nil
}

func (node *CommentWriter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		if pkt.Comment != "" {
			node.file.WriteString(pkt.Comment + "\n")
		}
		for _, output := range node.outputs {
			output <- pkt
		}
	}
}

func (node *CommentWriter) closeDown() {
	node.file.Close()
	log.Printf("closing down CommentWriter to file %s", node.file.Name())
	for _, output := range node.outputs {
		close(output)
	}
}
