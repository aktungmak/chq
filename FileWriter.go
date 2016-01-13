package main

import (
	"bufio"
	"log"
	"os"
)

// A FileWriter will write all packets it
// receives to a file.
// It passes through TS packets unmodified.
type FileWriter struct {
	file     *os.File
	writer   *bufio.Writer
	FileName string
	TsNode
}

//register with global AvailableNodes map
func init() {
	AvailableNodes.Register("FileWriter", NewFileWriter)
}

func NewFileWriter(fname string) (*FileWriter, error) {
	// try to open file
	fh, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}

	node := &FileWriter{}
	node.file = fh
	node.writer = bufio.NewWriterSize(fh, TS_PKT_SIZE*FILE_CHUNK_SIZE)
	node.input = make(chan TsPacket, CHAN_BUF_SIZE)
	node.FileName = fname

	go node.process()
	return node, nil
}

func (node *FileWriter) process() {
	defer node.closeDown()
	for pkt := range node.input {
		node.PktsIn++
		node.Send(pkt)
		node.writer.Write(pkt.bytes)
	}
}

func (node *FileWriter) closeDown() {
	node.Active = false
	node.writer.Flush()
	node.file.Close()
	log.Printf("closing down FileWriter to file %s", node.file.Name())
	node.output.Close()
}
