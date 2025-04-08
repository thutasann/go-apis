package main

import (
	"bytes"
	"testing"
)

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: DefaultPathTransformFunc,
	}
	s := NewStore(opts)

	data := bytes.NewReader([]byte("some jpb bytes"))
	if err := s.WriteStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}
}
