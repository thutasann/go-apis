package main

import "sync"

// Key and Value struct
type KV struct {
	mu   sync.RWMutex      // mutex lock
	data map[string][]byte // Key Value data map
}

// Initialize New Key Value
func NewKV() *KV {
	return &KV{
		data: map[string][]byte{},
	}
}

// Set the Key and Value
func (kv *KV) Set(key, val string) error {
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.data[key] = []byte(val)
	return nil
}

// Get the value with key
func (kv *KV) Get(key string) ([]byte, bool) {
	kv.mu.Lock()
	defer kv.mu.RUnlock()
	val, ok := kv.data[key]
	return val, ok
}
