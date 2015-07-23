package main

import (
	"io"
	"log"
	"os"
)

type FileInput struct {
	file *os.File
	TsNode
}

func NewFileInput(fname string) (*FileInput, error) {
	var err error

	// try to open file
	fh, err := os.Open(fname)
	Check(err)

	f := &FileInput{}
	f.file = fh
	f.input = nil
	f.outputs = make([]chan<- TsPacket, 0)

	go f.process()
	return f, nil
}

func (node *FileInput) process() {
	defer node.closeDown()
	buf := make([]byte, TS_PKT_SIZE*FILE_CHUNK_SIZE)
	for {
		// check for sync
		n, err := node.file.Read(buf)
		if err == io.EOF {
			break
		}
		Check(err)
		if n < TS_PKT_SIZE {
			log.Printf("Couldn't get a full packet, only %d bytes", n)
			continue
		}
		for i := 0; i < n; i += TS_PKT_SIZE {
			if buf[i] != 0x47 {
				log.Print("no lock yet")
				node.file.Seek(-int64(n-i-1), 1)
				break
			}
			for _, output := range node.outputs {
				output <- NewTsPacket(buf[i : i+TS_PKT_SIZE])
			}
		}
	}
}

func (node *FileInput) closeDown() {
	log.Printf("closing down file input for file %s", node.file.Name())
	node.file.Close()
	for _, output := range node.outputs {
		close(output)
	}
}
