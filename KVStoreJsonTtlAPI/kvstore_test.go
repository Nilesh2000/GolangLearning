package main

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSetAndGet(t *testing.T) {
	t.Run("get returns value after set", func(t *testing.T) {
		store := NewKVStore()
		k, v := "key1", "value1"

		store.Set(k, v)

		val, ok := store.Get(k)
		if !ok {
			t.Errorf("%q must exist", k)
		}
		if v != val {
			t.Errorf("got %q want %q", val, v)
		}
	})

	t.Run("check for non-existent key", func(t *testing.T) {
		store := NewKVStore()
		k := "key1"

		_, ok := store.Get(k)
		if ok {
			t.Errorf("%q must not exist", k)
		}
	})

	t.Run("set overwrites existing key", func(t *testing.T) {
		store := NewKVStore()

		store.Set("dup", "1")
		store.Set("dup", "2")

		val, ok := store.Get("dup")
		if !ok {
			t.Errorf("expected key to exist")
		}
		if val != "2" {
			t.Errorf("got %q want '2'", val)
		}
	})
}

func TestDelete(t *testing.T) {
	store := NewKVStore()
	k, v := "key1", "value1"

	store.Set(k, v)
	store.Delete(k)

	_, ok := store.Get(k)
	if ok {
		t.Errorf("%q must not exist", k)
	}
}

func TestKeys(t *testing.T) {
	store := NewKVStore()

	store.Set("key1", "value1")
	store.Set("key2", "value2")
	store.Set("key3", "value3")

	got := store.Keys()
	want := []string{"key1", "key2", "key3"}

	if len(got) != len(want) {
		t.Errorf("got %q want %q", got, want)
	}
}

func TestConcurrentAccess(t *testing.T) {
	store := NewKVStore()
	var wg sync.WaitGroup

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			k := "key" + strconv.Itoa(i)
			v := "val" + strconv.Itoa(i)
			store.Set(k, v)
		}(i)
	}

	wg.Wait()

	k, v := "key500", "val500"
	val, ok := store.Get("key500")
	if !ok || val != "val500" {
		t.Errorf("got %q want %q", k, v)
	}

	keys := store.Keys()
	if len(keys) != 1000 {
		t.Errorf("got %d, expected 1000", len(keys))
	}
}

func TestSetWithTTL(t *testing.T) {
	store := NewKVStore()
	defer store.Stop()

	tests := []struct {
		name       string
		key        string
		value      string
		ttl        time.Duration
		sleep      time.Duration
		wantExists bool
	}{
		{
			name:       "Key exists before TTL expires",
			key:        "key1",
			value:      "value1",
			ttl:        500 * time.Millisecond,
			sleep:      300 * time.Millisecond,
			wantExists: true,
		},
		{
			name:       "Key expires after TTl",
			key:        "key2",
			value:      "value2",
			ttl:        200 * time.Millisecond,
			sleep:      300 * time.Millisecond,
			wantExists: false,
		},
		{
			name:       "Zero TTL means no expiration",
			key:        "key3",
			value:      "value3",
			ttl:        0,
			sleep:      100 * time.Millisecond,
			wantExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			store.SetWithTTL(tt.key, tt.value, tt.ttl)

			time.Sleep(tt.sleep)

			_, exists := store.Get(tt.key)
			if exists != tt.wantExists {
				t.Errorf("got exists=%v, want=%v", exists, tt.wantExists)
			}
		})
	}
}

func TestConcurrentSetWithTTL(t *testing.T) {
	store := NewKVStore()
	defer store.Stop()

	var wg sync.WaitGroup
	wg.Add(100)
	for i := 1; i <= 100; i++ {
		go func() {
			defer wg.Done()
			key := "key" + strconv.Itoa(i)
			store.SetWithTTL(key, "value", 100*time.Millisecond)
		}()
	}

	wg.Wait()

	for i := 1; i <= 100; i++ {
		key := "key" + strconv.Itoa(i)
		if _, exists := store.Get(key); !exists {
			t.Errorf("key %s should exist", key)
		}
	}

	// Wait for expiration
	time.Sleep(150 * time.Millisecond)

	for i := 1; i <= 100; i++ {
		key := "key" + strconv.Itoa(i)
		if _, exists := store.Get(key); exists {
			t.Errorf("key %s should have expired", key)
		}
	}

}

func TestCleanupExpired(t *testing.T) {
	store := NewKVStore()
	defer store.Stop()

	for i := 1; i <= 10; i++ {
		key := "key" + strconv.Itoa(i)
		val := "value" + strconv.Itoa(i)
		store.SetWithTTL(key, val, 50*time.Millisecond)
	}

	time.Sleep(200 * time.Millisecond)

	if len(store.dict) != 0 {
		t.Errorf("store should be empty, got %d items", len(store.dict))
	}
}
