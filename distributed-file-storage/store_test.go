package main

import (
	"bytes"
	"fmt"
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

	data := bytes.NewReader([]byte("some jpb bytes"))
	if err := s.WriteStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}
}
