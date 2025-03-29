package main

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type item struct {
	Value      string
	Expiration int64
}

type KVStore struct {
	dict         map[string]*item
	mu           sync.RWMutex
	snapshotFile string
	saveInterval time.Duration
	stop         chan struct{}
}

func NewKVStore() *KVStore {
	store := &KVStore{
		dict:         make(map[string]*item),
		stop:         make(chan struct{}),
		snapshotFile: "", // Disabled by default
		saveInterval: 0,  // Disabled by default
	}

	go store.cleanupExpired()
	return store
}

func (s *KVStore) WithSnapshotFile(filename string) *KVStore {
	s.snapshotFile = filename
	return s
}

func (s *KVStore) WithSaveInterval(interval time.Duration) *KVStore {
	s.saveInterval = interval
	return s
}

func (s *KVStore) Initialize() *KVStore {
	if s.snapshotFile != "" {
		if data, err := os.ReadFile(s.snapshotFile); err == nil {
			json.Unmarshal(data, &s.dict)
		}
	}

	if s.snapshotFile != "" && s.saveInterval > 0 {
		go s.periodicSave()
	}
	return s
}

func (s *KVStore) Set(key string, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.dict[key] = &item{Value: value}
}

func (s *KVStore) SetWithTTL(key string, value string, ttl time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var exp int64
	if ttl > 0 {
		exp = time.Now().Add(ttl).UnixNano()
	}

	s.dict[key] = &item{
		Value:      value,
		Expiration: exp,
	}
}

func (s *KVStore) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	item, exists := s.dict[key]
	if !exists {
		return "", false
	}

	if item.Expiration > 0 && time.Now().UnixNano() > item.Expiration {
		return "", false
	}

	return item.Value, true
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
		if v.Expiration == 0 || now <= v.Expiration {
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
				if v.Expiration > 0 && now > v.Expiration {
					delete(s.dict, k)
				}
			}
			s.mu.Unlock()
		case <-s.stop:
			return
		}
	}
}

func (s *KVStore) periodicSave() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.saveToDisk()
		case <-s.stop:
			s.saveToDisk()
			return
		}
	}
}

func (s *KVStore) saveToDisk() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	data, err := json.Marshal(s.dict)
	if err != nil {
		log.Println("failed to marshal store:", err)
		return
	}
	if err := os.WriteFile(s.snapshotFile, data, 0644); err != nil {
		log.Println("failed to write snapshot:", err)
	}
}

func (s *KVStore) Stop() {
	close(s.stop)
}
