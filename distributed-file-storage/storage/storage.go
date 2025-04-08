/*
- A Content Addressable Storage
*/
package storage

import (
	"io"
	"log"
	"os"
)

// Path Transform Function
type PathTransformFunc func(string) string

// Store Options Struct
type StoreOpts struct {
	PathTransformFunc PathTransformFunc // Path Transform Func
}

// Store Struct
type Store struct {
	StoreOpts
}

// Get Default Path Transform Function
var DefaultPathTransformFunc = func(key string) string {
	return key
}

// Intitailize New Store
func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

// Write Stream
func (s *Store) WriteStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)

	if err := os.MkdirAll(pathName, os.ModePerm); err != nil {
		return err
	}

	filename := "somefilename"
	pathAndFilename := pathName + "/" + filename

	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathAndFilename)

	return nil
}
