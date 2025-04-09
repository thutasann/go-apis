package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strings"
)

// Default Path Transform Function
var DefaultPathTransformFunc = func(key string) string {
	return key
}

// CASPathTransformFunc takes a string key and transforms it into a deterministic, nested directory path based on the SHA-1 hash of the key.
//
// It's typically used in content-addressable storage (CAS) systems to organize files into a hierarchical directory structure, avoiding too many files in a single directory.
func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))          // [20]byte => []byte => [:]
	hashStr := hex.EncodeToString(hash[:]) // Convert the hash to a 40-character hex string

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Original: hashStr,
	}
}

// Path Transform Function
type PathTransformFunc func(string) PathKey

// Path Key
type PathKey struct {
	Pathname string // Path Name
	Original string // Original
}

// Store Options Struct
type StoreOpts struct {
	PathTransformFunc PathTransformFunc // Path Transform Func
}

// Store Struct
type Store struct {
	StoreOpts
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
	if err := os.MkdirAll(pathName.Pathname, os.ModePerm); err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	io.Copy(buf, r)

	filenameBytes := md5.Sum(buf.Bytes())
	filename := hex.EncodeToString(filenameBytes[:])
	pathAndFilename := pathName.Pathname + "/" + filename

	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathAndFilename)

	return nil
}
