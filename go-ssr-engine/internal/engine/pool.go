package engine

import (
	"bytes"
	"sync"
)

// bufferPool reduces GC pressure by reusing buffers.
//
// Rules:
// - Always Reset() before reuse
// - Never store excessively large buffers back
// - Keep pool usage in hot path only
var bufferPool = sync.Pool{
	New: func() any {
		return new(bytes.Buffer)
	},
}

// getBuffer retrieves reusable
func getBuffer() *bytes.Buffer {
	buf := bufferPool.Get().(*bytes.Buffer)
	buf.Reset()
	return buf
}

// putBuffer returns buffer to pool.
// Prevent unbounded memory retention by dropping huge buffers.
func putBuffer(buf *bytes.Buffer) {
	if buf.Cap() > 64*1024 {
		// Drop overlay large buffers to avoid memory bloat.
		return
	}
	bufferPool.Put(buf)
}
