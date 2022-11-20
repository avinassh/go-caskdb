package caskdb

import (
	"os"
	"testing"
)

func TestDiskStore_Get(t *testing.T) {
	store := NewDiskStore("test.db")
	defer os.Remove("test.db")
	store.Set("name", "jojo")
	if val := store.Get("name"); val != "jojo" {
		t.Errorf("Get() = %v, want %v", val, "jojo")
	}
}

func TestDiskStore_GetInvalid(t *testing.T) {
	store := NewDiskStore("test.db")
	defer os.Remove("test.db")
	if val := store.Get("some key"); val != "" {
		t.Errorf("Get() = %v, want %v", val, "")
	}
}
