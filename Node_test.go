package main

import (
	"testing"
)

func BenchmarkSend(b *testing.B) {
	node := TsNode{}
	out1 := make(chan TsPacket)
	out2 := make(chan TsPacket)

	node.RegisterListener(out1)
	node.RegisterListener(out2)

	pkt := TsPacket{}
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
