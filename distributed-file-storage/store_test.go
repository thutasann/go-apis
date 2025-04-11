package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "superherosbestpics"
	pathKey := CASPathTransformFunc(key)
	expectedPathName := "9b865/0f1f3/64b62/fa398/fbecc/2b99e/0420a/9e169"
	expectedOriginalName := "9b8650f1f364b62fa398fbecc2b99e0420a9e169"

	if pathKey.PathName != expectedPathName {
		t.Error(t, "have %s want %s", pathKey, expectedPathName)
	}

	if pathKey.FileName != expectedOriginalName {
		t.Error(t, "have %s want %s", pathKey, expectedPathName)
	}
	fmt.Println(pathKey)
}

func TestStore(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momspecials"
	data := []byte("some jpb bytes")

	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, _ := io.ReadAll(r)

	fmt.Println(string(b))

	if string(b) != string(data) {
		t.Errorf("want %s have %s", data, b)
	}

	s.Delete(key)
}

func TestDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momspecials"
	data := []byte("some jpb bytes")

	if err := s.WriteStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}
