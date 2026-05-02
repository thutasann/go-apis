package main

import (
	"os"
	"sync"
)

type input struct {
	mu          sync.Mutex
	pressedKey  byte
	frameKey    byte
}

func newInput() *input {
	i := &input{}
	go i.readLoop()
	return i
}

func (i *input) readLoop() {
	b := make([]byte, 1)
	for {
		if _, err := os.Stdin.Read(b); err != nil {
			return
		}
		i.mu.Lock()
		i.frameKey = b[0]
		i.pressedKey = b[0]
		i.mu.Unlock()
	}
}

func (i *input) update() {
	i.mu.Lock()
	defer i.mu.Unlock()
}

func (i *input) consumeFrameKey() byte {
	i.mu.Lock()
	defer i.mu.Unlock()
	key := i.frameKey
	i.frameKey = 0
	return key
}

func (i *input) key() byte {
	i.mu.Lock()
	defer i.mu.Unlock()
	return i.pressedKey
}
