package main

import (
	"sync"
	"testing"
	"time"
)

var (
	N       = 3
	testPkt = TsPacket{}
	timeout = time.Second
)

type ListenFunc func(int, *Broadcaster, *sync.WaitGroup)

func setupN(f ListenFunc) (*Broadcaster, *sync.WaitGroup) {
	var b Broadcaster
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go f(i, &b, &wg)
	}
	wg.Wait()
	return &b, &wg
}

func TestRegisterChan(t *testing.T) {
	b, wg := setupN(func(i int, b *Broadcaster, wg *sync.WaitGroup) {
		ch := make(chan TsPacket)
		b.RegisterChan(ch)
		wg.Done()
		select {
		case v := <-ch:
			if v.Comment != testPkt.Comment {
				t.Error("bad value received")
			}
		case <-time.After(timeout):
			t.Error("receive timed out")
		}
		wg.Done()
	})
	wg.Add(N)
	b.Send(testPkt)
	wg.Wait()

}
