package webhook

import (
	"fmt"

	consul "github.com/hashicorp/consul/api"
)

type consulStore struct {
	kv *consul.KV
}

func (store *consulStore) Put(key string, val []byte) error {
	_, err := store.kv.Put(&consul.KVPair{
		Key:   key,
		Value: val,
	}, nil)
	return err
}

func (store *consulStore) Get(key string) ([]byte, error) {
	kv, _, err := store.kv.Get(key, nil)
	if err != nil {
		return nil, err
	}
	if kv == nil {
		return nil, fmt.Errorf("consulStoreError: key %s is not found", key)
	}

	return kv.Value, nil
}

func (store *consulStore) Delete(key string) error {
	_, err := store.kv.Delete(key, nil)
	return err
}

func (store *consulStore) List(key string) ([]string, error) {
	var vals []string
	kvPairs, _, err := store.kv.List(key, nil)

	if err != nil {
		return nil, err
	}

	for _, kv := range kvPairs {
		vals = append(vals, string(kv.Value))
	}
	return vals, err
}

// NewConsulStore returns a new consul store
func NewConsulStore() Store {
	client, err := consul.NewClient(consul.DefaultConfig())
	if err != nil {
		panic(err)
	}

	kv := client.KV()
	return &consulStore{
		kv: kv,
	}
}
