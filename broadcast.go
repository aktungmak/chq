package main

import "sync"

// Broadcaster implements a broadcast channel.
// The zero value is a usable unbuffered channel.
type Broadcaster struct {
	m         sync.Mutex
	listeners map[int]chan<- TsPacket // lazy init
	nextId    int
	capacity  int
	closed    bool
}

// New returns a new Broadcaster with the given capacity (0 means unbuffered).
func New(n int) *Broadcaster {
	return &Broadcaster{capacity: n}
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
		b.listeners = make(map[int]chan<- TsPacket)
	}
	for b.listeners[b.nextId] != nil {
		b.nextId++
	}
	if b.closed {
		close(ch)
	}
	b.listeners[b.nextId] = ch
}
func (b *Broadcaster) UnRegisterChan(ch chan TsPacket) {
	b.m.Lock()
	defer b.m.Unlock()
	if b.listeners == nil {
		b.listeners = make(map[int]chan<- TsPacket)
	}
	for id, ch := range b.listeners {
		if ch == toremove {
			delete(b.listeners, id)
			break
		}
	}
}
