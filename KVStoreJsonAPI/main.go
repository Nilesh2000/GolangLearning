package main

import (
	"log"
	"net/http"
	"sync"
)

type KVStore struct {
	dict map[string]string
	mu   sync.Mutex
}

func NewKVStore() *KVStore {
	return &KVStore{
		dict: make(map[string]string),
	}
}

func (s *KVStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.dict[key] = value
}

func (s *KVStore) Get(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	val, exists := s.dict[key]
	if !exists {
		return "", false
	}
	return val, true
}

func (s *KVStore) Delete(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.dict, key)
}

// Keys returns a slice of all keys in the store.
// The order of keys is not guaranteed.
func (s *KVStore) Keys() []string {
	s.mu.Lock()
	defer s.mu.Unlock()

	keys := make([]string, 0, len(s.dict))
	for k := range s.dict {
		keys = append(keys, k)
	}
	return keys
}

func main() {
	store := NewKVStore()
	s := NewServer(store)
	log.Fatal(http.ListenAndServe(":8080", s.mux))
}
