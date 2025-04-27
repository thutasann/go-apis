package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Default Root Folder name
const defaultRootFolderName = "thuta_network"

// Path Key
type PathKey struct {
	PathName string // PathKey's Path Name
	FileName string // PathKey's FileName
}

// Path Transform Function
type PathTransformFunc func(string) PathKey

// Store Options Struct
type StoreOpts struct {
	Root              string            // Root is the folder name of the root, containing all the folders/files of the system
	PathTransformFunc PathTransformFunc // Path Transform Func
}

// Store Struct
type Store struct {
	StoreOpts
}

// Default Path Transform Function (PathTransformFunc)
var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		PathName: key,
		FileName: key,
	}
}

// CASPathTransformFunc (PathTransformFunc) takes a string key and transforms it into a deterministic, nested directory path based on the SHA-1 hash of the key.
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
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

// Read Data
// Instead of copying directly to a reader, we first copy this into
// a buffer. Maybe just return the file from teh readstream
func (s *Store) Read(id string, key string) (int64, io.Reader, error) {
	return s.readStream(id, key)
}

// Read Stream
func (s *Store) readStream(id string, key string) (int64, io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	file, err := os.Open(fullPathWithRoot)
	if err != nil {
		return 0, nil, err
	}

	fi, err := file.Stat()
	if err != nil {
		return 0, nil, err
	}

	return fi.Size(), file, nil
}

// Write Data to the Disk
func (s *Store) Write(id string, key string, r io.Reader) (int64, error) {
	return s.writeStream(id, key, r)
}

// Write Stream
func (s *Store) writeStream(id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWrtiing(id, key)
	if err != nil {
		return 0, err
	}
	return io.Copy(f, r)
}

// Write Decrypt
func (s *Store) WriteDecrypt(encKey []byte, id string, key string, r io.Reader) (int64, error) {
	f, err := s.openFileForWrtiing(id, key)
	if err != nil {
		return 0, err
	}
	n, _ := CopyDecrypt(encKey, r, f)
	return int64(n), nil
}

// Open file for writing
func (s *Store) openFileForWrtiing(id string, key string) (*os.File, error) {
	pathKey := s.PathTransformFunc(key)
	pathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.PathName)
	if err := os.MkdirAll(pathNameWithRoot, os.ModePerm); err != nil {
		return nil, err
	}

	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())

	return os.Create(fullPathWithRoot)
}

// Delete Data
func (s *Store) Delete(id string, key string) error {
	pathKey := s.PathTransformFunc(key)

	defer func() {
		log.Printf("deleted [%s] from disk", pathKey.FileName)
	}()

	firstPathNameWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FirstPathName())
	return os.RemoveAll(firstPathNameWithRoot)
}

// Clear the Root
func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

// Check Has Path
func (s *Store) Has(id string, key string) bool {
	pathKey := s.PathTransformFunc(key)
	fullPathWithRoot := fmt.Sprintf("%s/%s/%s", s.Root, id, pathKey.FullPath())
	_, err := os.Stat(fullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}
