package main

import (
	"sync"
	"time"
)

type item struct {
	value      string
	expiration int64
}

type KVStore struct {
	dict map[string]*item
	mu   sync.RWMutex
	stop chan struct{}
}

func NewKVStore() *KVStore {
	store := &KVStore{
		dict: make(map[string]*item),
		stop: make(chan struct{}),
	}
	go store.cleanupExpired()
	return store
}

func (s *KVStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.dict[key] = &item{value: value, expiration: 0}
}

func (s *KVStore) SetWithTTL(key string, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	s.dict[key] = &item{
		value:      value,
		expiration: exp,
	}
}

func (s *KVStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.dict[key]
	if !exists {
		return "", false
	}

	if item.expiration > 0 && time.Now().UnixNano() > item.expiration {
		return "", false
	}

	return item.value, true
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
	now := time.Now().UnixNano()
	for k, v := range s.dict {
		if v.expiration == 0 || now <= v.expiration {
			keys = append(keys, k)
		}
	}
	return keys
}

func (s *KVStore) cleanupExpired() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.mu.Lock()
			now := time.Now().UnixNano()
			for k, v := range s.dict {
				if v.expiration > 0 && now > v.expiration {
					delete(s.dict, k)
				}
			}
			s.mu.Unlock()
		case <-s.stop:
			return
		}
	}
}

func (s *KVStore) Stop() {
	close(s.stop)
}
