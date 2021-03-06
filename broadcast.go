package main

import (
	"sync"
)

// Broadcaster implements a broadcast channel.
// The zero value is a usable channel.
type Broadcaster struct {
	m         sync.Mutex
	listeners map[int]chan TsPacket // lazy init
	nextId    int
	closed    bool
}

// Send broadcasts a message to the channel.
// If the channel is closed, do nothing.
func (b *Broadcaster) Send(v TsPacket) {
	b.m.Lock()
	defer b.m.Unlock()
	if b.closed {
		return
	}
	for _, l := range b.listeners {
		l <- v
	}
}

// Close closes the channel, disabling the sending of further messages.
func (b *Broadcaster) Close() {
	b.m.Lock()
	defer b.m.Unlock()
	b.closed = true
	for _, l := range b.listeners {
		close(l)
	}
}

func (b *Broadcaster) RegisterChan(ch chan TsPacket) {
	b.m.Lock()
	defer b.m.Unlock()
	if b.listeners == nil {
		b.listeners = make(map[int]chan TsPacket)
	}
	for b.listeners[b.nextId] != nil {
		b.nextId++
	}
	if b.closed {
		// there is no way to know if ch is already closed
		// so we have to let it panic and then recover
		defer func() {
			if r := recover(); r != nil {
				b.listeners[b.nextId] = ch
			}
		}()
		close(ch)
	}
	b.listeners[b.nextId] = ch
}
func (b *Broadcaster) UnRegisterChan(ch chan TsPacket) {
	b.m.Lock()
	defer b.m.Unlock()
	if b.listeners == nil {
		b.listeners = make(map[int]chan TsPacket)
	}
	for i, c := range b.listeners {
		if c == ch {
			delete(b.listeners, i)
			break
		}
	}
}

func (b *Broadcaster) GetOutputs() []chan TsPacket {
	b.m.Lock()
	defer b.m.Unlock()
	ret := make([]chan TsPacket, 0, len(b.listeners))
	for _, ch := range b.listeners {
		ret = append(ret, ch)
	}
	return ret
}
