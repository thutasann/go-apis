package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"strings"
)

// Path Key
type PathKey struct {
	PathName string // PathKey's Path Name
	FileName string // PathKey's FileName
}

// Path Transform Function
type PathTransformFunc func(string) PathKey

// Store Options Struct
type StoreOpts struct {
	PathTransformFunc PathTransformFunc // Path Transform Func
}

// Store Struct
type Store struct {
	StoreOpts
}

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
		PathName: strings.Join(paths, "/"),
		FileName: hashStr,
	}
}

// Get First Path Name from the PathKey
func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.PathName, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

// Get Full Path From PathKey (PathName and FileName)
func (p PathKey) FullPath() string {
	return fmt.Sprintf("%s/%s", p.PathName, p.FileName)
}

// Intitailize New Store
func NewStore(opts StoreOpts) *Store {
	return &Store{
		StoreOpts: opts,
	}
}

// Read Stream
func (s *Store) ReadStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	return os.Open(pathKey.FullPath())
}

// Write Stream
func (s *Store) WriteStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	if err := os.MkdirAll(pathKey.PathName, os.ModePerm); err != nil {
		return err
	}

	fullPath := pathKey.FullPath()

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, fullPath)

	return nil
}

// Read Data
func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.ReadStream(key)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)

	return buf, err
}

// Delete Data
func (s *Store) Delete(key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	// if err := os.RemoveAll(pathKey.FullPath()); err != nil {
	// 	return err
	// }

	return os.RemoveAll(pathKey.FirstPathName())
}

// Check Has Path
func (s *Store) Has(key string) bool {
	pathKey := s.PathTransformFunc(key)
	_, err := os.Stat(pathKey.FullPath())
	return !errors.Is(err, fs.ErrNotExist)
}
