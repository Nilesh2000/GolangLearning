package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestPersistence_ReloadOnRestart(t *testing.T) {
	snapshotFile := "store.snapshot.json"
	os.Remove(snapshotFile)

	store := NewKVStore().
		WithSnapshotFile(snapshotFile).
		WithSaveInterval(5 * time.Second).
		Initialize()

	store.Set("key1", "value1")
	store.Set("key2", "value2")
	store.Stop()

	time.Sleep(time.Second) // Some breathing space

	reloaded := NewKVStore().
		WithSnapshotFile(snapshotFile).
		Initialize()
	defer reloaded.Stop()

	if val, _ := reloaded.Get("key1"); val != "value1" {
		t.Errorf("expected key1 to be restored from snapshot")
	}

	if val, _ := reloaded.Get("key2"); val != "value2" {
		t.Errorf("expected key2 to be restored from snapshot")
	}
}

func TestPersistence_RespectsTTL(t *testing.T) {
	snapshotFile := "store.snapshot.json"
	os.Remove(snapshotFile)

	store := NewKVStore().
		WithSnapshotFile(snapshotFile).
		WithSaveInterval(5 * time.Second).
		Initialize()

	store.SetWithTTL("key1", "value1", time.Second)
	store.SetWithTTL("key2", "value2", time.Minute)
	time.Sleep(6 * time.Second) // Wait for Auto-Save
	store.Stop()

	reloaded := NewKVStore().
		WithSnapshotFile(snapshotFile).
		Initialize()
	defer reloaded.Stop()

	if _, exists := reloaded.Get("key1"); exists {
		t.Errorf("expected expired key1 to be gone")
	}

	if val, _ := reloaded.Get("key2"); val != "value2" {
		t.Errorf("expected key2 to be restored from snapshot")
	}
}

func TestPersistence_FileWrittenCorrectly(t *testing.T) {
	snapshotFile := "store.snapshot.json"
	os.Remove(snapshotFile)

	store := NewKVStore().
		WithSnapshotFile(snapshotFile).
		WithSaveInterval(5 * time.Second).
		Initialize()

	store.Set("a", "1")
	store.Set("b", "2")
	store.Stop()

	time.Sleep(time.Second) // Some breathing space

	raw, err := os.ReadFile(snapshotFile)
	if err != nil {
		t.Fatalf("expected file to exist")
	}

	var decoded map[string]map[string]interface{}
	err = json.Unmarshal(raw, &decoded)
	if err != nil {
		t.Errorf("invalid snapshot format")
	}

	if len(decoded) != 2 {
		t.Fatalf("unexpected snapshot content")
	}
}
