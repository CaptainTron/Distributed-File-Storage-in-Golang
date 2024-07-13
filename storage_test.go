package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)

func newStore() *Store {
	opts := StoreOpts{
		PathTransformFunc: CASPathTransformFunc,
	}
	return NewStore(opts)
}

func teardown(t *testing.T, s *Store) {
	if err := s.Clear(); err != nil {
		t.Error(err)
	}
}

func TestPathTransformFunc(t *testing.T) {
	key := "mombestpicture"
	pathKey := CASPathTransformFunc(key)

	expectedFilename := "cf5d4b01c4d9438c22c56c832f83bd3e8c6304f9"
	expectedPathname := "cf5d4/b01c4/d9438/c22c5/6c832/f83bd/3e8c6/304f9"

	if pathKey.Filename != expectedFilename {
		t.Errorf("have %s want %s", pathKey.Filename, expectedFilename)
	}

	if pathKey.Pathname != expectedPathname {
		t.Errorf("have %s want %s", pathKey.Pathname, expectedPathname)
	}
}

func TestStore(t *testing.T) {
	s := newStore()
	//defer teardown(t, s);

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpg bytes")
		if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		// Read the content of file in folders and subfolders
		r, err := s.Read(key)
		if err != nil {
			t.Error(err)
		}

		// Print the content of file on the screen
		b, _ := ioutil.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("want %s have %s", data, b)
		}
		fmt.Println(string(b))

		// Delete the file and its folders
		if err := s.Delete(key); err != nil {
			t.Error(err)
		}

		// Check it exists or not
		if ok := s.Has(key); !ok {
			t.Errorf("expected to have key %s", key)
		}
	}
}
