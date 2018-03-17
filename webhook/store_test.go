package webhook

import (
	"testing"
)

func TestStore_GetNonExistingKeyInMemory(t *testing.T) {
	store := NewStore(InMemory)
	key := "hello"

	_, err := store.Get(key)
	if err == nil {
		t.Errorf("expected error, got none %v", err)
	}
}

func TestStore_PutGetInMemory(t *testing.T) {
	store := NewStore(InMemory)
	key := "hello"
	val := []byte("world")

	err := store.Put(key, val)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	gotVal, err := store.Get(key)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if string(gotVal) != string(val) {
		t.Errorf("expected %v, got %v", val, gotVal)
	}
}

func TestStore_GetNonExistingConsul(t *testing.T) {
	store := NewStore(Consul)
	key := "key/non/existing"
	_, err := store.Get(key)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}
}
func TestStore_PutGetConsul(t *testing.T) {
	store := NewStore(Consul)
	key := "hello"
	val := []byte("world")

	err := store.Put(key, val)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	gotVal, err := store.Get(key)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if string(gotVal) != string(val) {
		t.Errorf("expected %v, got %v", val, gotVal)
	}
}

// func TestStore_PutPrefixConsul(t *testing.T) {

// }
