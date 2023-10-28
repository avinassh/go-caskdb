package caskdb

import (
	"os"
	"testing"
)

func TestDiskStore_Get(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")
	store.Set("name", "jojo")
	if val := store.Get("name"); val != "jojo" {
		t.Errorf("Get() = %v, want %v", val, "jojo")
	}
}

func TestDiskStore_GetInvalid(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")
	if val := store.Get("some key"); val != "" {
		t.Errorf("Get() = %v, want %v", val, "")
	}
}

func TestDiskStore_SetWithPersistence(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")

	tests := map[string]string{
		"crime and punishment": "dostoevsky",
		"anna karenina":        "tolstoy",
		"war and peace":        "tolstoy",
		"hamlet":               "shakespeare",
		"othello":              "shakespeare",
		"brave new world":      "huxley",
		"dune":                 "frank herbert",
	}
	for key, val := range tests {
		store.Set(key, val)
		if store.Get(key) != val {
			t.Errorf("Get() = %v, want %v", store.Get(key), val)
		}
	}
	store.Close()
	store, err = NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	for key, val := range tests {
		if store.Get(key) != val {
			t.Errorf("Get() = %v, want %v", store.Get(key), val)
		}
	}
	store.Close()
}

func TestDiskStore_Delete(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")

	tests := map[string]string{
		"crime and punishment": "dostoevsky",
		"anna karenina":        "tolstoy",
		"war and peace":        "tolstoy",
		"hamlet":               "shakespeare",
		"othello":              "shakespeare",
		"brave new world":      "huxley",
		"dune":                 "frank herbert",
	}
	for key, val := range tests {
		store.Set(key, val)
	}

	// only for tests
	deletedKeys := []string{"hamlet", "dune", "othello"}
	//delete few keys
	for _, k := range deletedKeys {
		store.Delete(k)
	}
	store.Close()

	store, err = NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}

	//check for deletion
	for _, dkeys := range deletedKeys {
		actualVal := store.Get(dkeys)
		if actualVal != TombStoneVal {
			t.Errorf("Get() = %s, want %s", actualVal, TombStoneVal)
		}
	}
	store.Close()
}
