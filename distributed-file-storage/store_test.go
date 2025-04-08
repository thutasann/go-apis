package main

import (
	"bytes"
	"fmt"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "superherosbestpics"
	pathname := CASPathTransformFunc(key)
	expectedPathName := "9b865/0f1f3/64b62/fa398/fbecc/2b99e/0420a/9e169"
	if pathname != expectedPathName {
		t.Error(t, "have %s want %s", pathname, expectedPathName)
	}
	fmt.Println(pathname)
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
