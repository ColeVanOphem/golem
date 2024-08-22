package main

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestPathTransformFunc(t *testing.T) {
	key := "mybestpicture"
	pathKey := CASPathTransformFunc(key)
	expectedOriginalName := "be17b32c2870b1c0c73b59949db6a3be7814dd23"
	expectedPath := "be17b/32c28/70b1c/0c73b/59949/db6a3/be781/4dd23"

	if pathKey.PathName != expectedPath {
		t.Errorf("Received path %s but expected path %s", pathKey.PathName, expectedPath)
	}
	if pathKey.FileName != expectedOriginalName {
		t.Errorf("Received original name %s but expected %s", pathKey.FileName, expectedOriginalName)
	}
}

// func TestDelete(t *testing.T) {
// 	opts := StoreOpts{
// 		PathTransformFunc: CASPathTransformFunc,
// 	}

// 	s := NewStore(opts)

// 	key := "mybestpicture"
// 	data := []byte("some jpg bytes")

// 	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
// 		t.Error(err)
// 	}

// 	if err := s.Delete(key); err != nil {
// 		t.Error(err)
// 	}
// }

func TestStore(t *testing.T) {
	s := newStore()
	id := generateID()
	defer teardown(t, s)

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf("foo_%d", i)
		data := []byte("some jpg bytes")

		if _, err := s.writeStream(id, key, bytes.NewReader(data)); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); !ok {
			t.Errorf("Expected to have key %s", key)
		}

		_, r, err := s.Read(id, key)
		if err != nil {
			t.Error(err)
		}

		b, _ := io.ReadAll(r)
		if string(b) != string(data) {
			t.Errorf("Expected %s but read %s", data, b)
		}

		if err := s.Delete(id, key); err != nil {
			t.Error(err)
		}

		if ok := s.Has(id, key); ok {
			t.Errorf("Expected to NOT have key %s", key)
		}
	}
}

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
