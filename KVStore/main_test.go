package main

import (
	"strconv"
	"sync"
	"testing"
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
