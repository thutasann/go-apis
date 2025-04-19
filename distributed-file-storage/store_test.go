package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

// Test Path Transform Function
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

// Test Store
func TestStore(t *testing.T) {
	s := newStore()
	defer teardown(t, s)

	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("food_%d", i)
		data := []byte("some jpb bytes")

		size, err := s.writeStream(key, bytes.NewReader(data))
		if err != nil {
			t.Error(err)
		}

		fmt.Println("size:", size)

		if ok := s.Has(key); !ok {
			t.Errorf("expected to have key %s", key)
		}

		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}

		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(key); ok {
			t.Errorf("expected to NOT have key %s", key)
		}
	}
}

// Test Delete
func TestDelete(t *testing.T) {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	key := "momspecials"
	data := []byte("some jpb bytes")

	size, err := s.writeStream(key, bytes.NewReader(data))
	if err != nil {
		t.Error(err)
	}

	fmt.Println("written size: ", size)

	if err := s.Delete(key); err != nil {
		t.Error(err)
	}
}

// Get New Store
func newStore() *Store {
	opts := StoreOpts{
		Root:              defaultRootFolderName,
		PathTransformFunc: CASPathTransformFunc,
	}
	s := NewStore(opts)
	return s
}

// Test Teardown
func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}
