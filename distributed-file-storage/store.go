package main

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

// CASPathTransformFunc takes a string key and transforms it into a deterministic, nested directory path based on the SHA-1 hash of the key.
//
// It's typically used in content-addressable storage (CAS) systems to organize files into a hierarchical directory structure, avoiding too many files in a single directory.
func CASPathTransformFunc(key string) string {
	hash := sha1.Sum([]byte(key))          // [20]byte => []byte => [:]
	hashStr := hex.EncodeToString(hash[:]) // Convert the hash to a 40-character hex string

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return strings.Join(paths, "/")
}

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
