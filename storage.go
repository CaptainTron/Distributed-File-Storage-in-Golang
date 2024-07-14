package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

const defaultRootFolderName = "vaibhav"

func CASPathTransformFunc(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blocksize := 5
	sliceLen := len(hashStr) / blocksize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blocksize, (i*blocksize)+blocksize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		Pathname: strings.Join(paths, "/"),
		Filename: hashStr,
	}
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) Fullpath() string {
	return fmt.Sprintf("%s/%s", p.Pathname, p.Filename)
}

type StoreOpts struct {
	// Root is the folder name of the root, containing all the folder/files of systems
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key,
	}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	// If transformFunc not defined
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

func (s *Store) Clear() error {
	return os.RemoveAll(s.Root)
}

// Returns first filepathname
func (p PathKey) FirstPathName() string {
	paths := strings.Split(p.Pathname, "/")
	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (s *Store) Has(key string) bool {
	PathKey := s.PathTransformFunc(key)
	FullPathWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.Fullpath())
	_, err := os.Stat(FullPathWithRoot)
	return !errors.Is(err, os.ErrNotExist)
}

// Delete the file contents and its children
func (s *Store) Delete(key string) error {
	PathKey := s.PathTransformFunc(key)
	defer func() {
		log.Printf("deleted [%s] from disk", PathKey.Filename)
	}()
	FirstPathNameWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.FirstPathName())
	return os.RemoveAll(FirstPathNameWithRoot)
}
func (s *Store) Write(key string, r io.Reader) error {
	return s.writeStream(key, r)
}

func (s *Store) Read(key string) (io.Reader, error) {
	f, err := s.readStream(key)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	return buf, err
}

func (s *Store) readStream(key string) (io.Reader, error) {
	PathKey := s.PathTransformFunc(key)
	PathKeyWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.Fullpath())
	return os.Open(PathKeyWithRoot)
}

func (s *Store) writeStream(key string, r io.Reader) error {
	PathKey := s.PathTransformFunc(key)
	pathnameWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.Pathname)
	if err := os.MkdirAll(pathnameWithRoot, os.ModePerm); err != nil {
		return err
	}

	pathandFilenameWithRoot := fmt.Sprintf("%s/%s", s.Root, PathKey.Fullpath())

	f, err := os.Create(pathandFilenameWithRoot)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathandFilenameWithRoot)
	return nil
}
