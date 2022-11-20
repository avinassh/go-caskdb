package caskdb

import "testing"

func TestMemoryStore_Get(t *testing.T) {
	store := NewMemoryStore()
	store.Set("name", "jojo")
	if val := store.Get("name"); val != "jojo" {
		t.Errorf("Get() = %v, want %v", val, "jojo")
	}
}

func TestMemoryStore_InvalidGet(t *testing.T) {
	store := NewMemoryStore()
	if val := store.Get("some rando key"); val != "" {
		t.Errorf("Get() = %v, want %v", val, "")
	}
}

func TestMemoryStore_Close(t *testing.T) {
	store := NewMemoryStore()
	if !store.Close() {
		t.Errorf("Close() failed")
	}
}
