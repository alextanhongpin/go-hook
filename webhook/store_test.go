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

func TestStore_PutPrefixConsul(t *testing.T) {
	store := NewStore(Consul)
	key := "key"
	val := "val"

	err := store.Put(key, []byte(val))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	err = store.Delete(key)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	_, err = store.Get(key)
	if err == nil {
		t.Errorf("expected error, got %v", err)
	}
}

func TestStore_ListConsul(t *testing.T) {
	store := NewStore(Consul)

	testTables := []struct {
		key, val string
	}{
		{"key/foo", "val1"},
		{"key/bar", "val2"},
	}

	for _, test := range testTables {
		err := store.Put(test.key, []byte(test.val))
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
	}

	vals, err := store.List("key")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	if len(vals) != 2 {
		t.Errorf("expected len to be 2, got %d", len(vals))
	}
}
