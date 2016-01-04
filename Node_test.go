package main

import (
	"testing"
)

var sampPkt = []byte{
	0x47, 0x84, 0x91, 0x33, 0x07, 0x10, 0x48, 0xC2, 0xD5, 0x83, 0xFE,
	0xE4, 0x0F, 0xD2, 0x4E, 0xDF, 0xD3, 0xF5, 0x80, 0xC1, 0x1A, 0xE5,
	0xE6, 0x26, 0x55, 0x09, 0xDB, 0x31, 0x40, 0x7F, 0xD1, 0x53, 0x26,
	0x7F, 0x80, 0x53, 0xE4, 0x56, 0x8B, 0x6F, 0x56, 0xE2, 0x83, 0x2B,
	0xFA, 0xEE, 0x3D, 0x88, 0x50, 0x3E, 0x85, 0x94, 0xE8, 0x46, 0x76,
	0xDD, 0xDD, 0x92, 0x69, 0xE5, 0xCE, 0xA9, 0x16, 0x1A, 0x8C, 0x77,
	0x8F, 0x02, 0x09, 0xA3, 0xB5, 0xB1, 0x65, 0xC6, 0xDE, 0x93, 0x95,
	0xEB, 0x29, 0x73, 0x1E, 0x0F, 0xAC, 0x81, 0xB8, 0x7E, 0xE8, 0x22,
	0xFD, 0x04, 0x86, 0xD2, 0xCB, 0x27, 0x10, 0x6E, 0xE8, 0x87, 0xBC,
	0x1B, 0xB4, 0xE2, 0x55, 0x9A, 0x66, 0xF1, 0x46, 0x44, 0x6C, 0xE8,
	0xE3, 0xA2, 0x34, 0xCA, 0x3C, 0x39, 0x24, 0xE6, 0x78, 0xF3, 0xA0,
	0xA7, 0xBA, 0x92, 0xA7, 0x9D, 0x17, 0xFB, 0x6C, 0x57, 0xDC, 0x83,
	0x7D, 0x97, 0x8A, 0xE2, 0x44, 0xE0, 0xAE, 0x5A, 0x88, 0x0A, 0xDE,
	0x02, 0x59, 0xF4, 0xBD, 0x85, 0xEB, 0x9A, 0x32, 0x0F, 0x3D, 0xA8,
	0xAF, 0xFF, 0xF5, 0x91, 0xCE, 0x7E, 0x87, 0x06, 0x6F, 0xE5, 0xCF,
	0xFE, 0x21, 0xB5, 0x3B, 0x78, 0x51, 0xB0, 0x6F, 0x48, 0xE7, 0xD6,
	0x8F, 0xF8, 0x67, 0x99, 0x2A, 0x5B, 0x56, 0x3F, 0xCF, 0xC7, 0xDA, 0x11}

func BenchmarkSend(b *testing.B) {
	node := TsNode{}
	out1 := make(chan TsPacket)
	out2 := make(chan TsPacket)

	node.RegisterListener(out1)
	node.RegisterListener(out2)

	pkt := NewTsPacket(sampPkt)
	go func() {
		for p := range out1 {
			p = p
		}
	}()
	go func() {
		for p := range out2 {
			p = p
		}
	}()

	node.Toggle()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.Send(pkt)
	}
	close(out1)
	close(out2)
}

func TestPidRemapper(t *testing.T) {
	node, _ := NewPidRemapper(1169, 202)
	ipkt := NewTsPacket(sampPkt)
	out := make(chan TsPacket)

	node.RegisterListener(out)
	node.Toggle()

	node.GetInputChan() <- ipkt
	opkt := <-out

	if opkt.Header.Pid != 202 {
		t.Errorf("opkt.Header.Pid is incorrect, got %d", opkt.Header.Pid)
	}
	if ((int(opkt.bytes[1])&31)<<8)+int(opkt.bytes[2]) != 202 {
		t.Errorf("pid in opkt.bytes is incorrect")
	}
}

func BenchmarkPidRemapper(b *testing.B) {
	node, _ := NewPidRemapper(1169, 202)
	out := make(chan TsPacket)

	node.RegisterListener(out)

	pkt := NewTsPacket(sampPkt)
	go func() {
		for p := range out {
			p = p
		}
	}()

	node.Toggle()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		node.GetInputChan() <- pkt
	}
	// close(out)
}
