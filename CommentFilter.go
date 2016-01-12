package main

import (
	"log"
)

// A CommentFilter will only pass TSPackets
// which have a non-empty Comment field.
// Concretely, len(Comment) > 0
type CommentFilter struct {
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("CommentFilter", NewCommentFilter)
}

func NewCommentFilter() (*CommentFilter, error) {
	node := &CommentFilter{}
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)

	go node.process()
	return node, nil
}

func (node *CommentFilter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		node.PktsIn++
		if len(pkt.Comment) > 0 {
			node.Send(pkt)
		}
	}
}

func (node *CommentFilter) closeDown() {
	node.Active = false
	log.Print("closing down CommentFilter")
	node.output.Close()
}
