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

func TestDiskStore_SetWithPersistence(t *testing.T) {
	store := NewDiskStore("test.db")
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
	store = NewDiskStore("test.db")
	for key, val := range tests {
		if store.Get(key) != val {
			t.Errorf("Get() = %v, want %v", store.Get(key), val)
		}
	}
	store.Close()
}

func TestDiskStore_Delete(t *testing.T) {
	store := NewDiskStore("test.db")
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
	for key, _ := range tests {
		store.Set(key, "")
	}
	store.Set("end", "yes")
	store.Close()

	store = NewDiskStore("test.db")
	for key := range tests {
		if store.Get(key) != "" {
			t.Errorf("Get() = %v, want '' (empty)", store.Get(key))
		}
	}
	if store.Get("end") != "yes" {
		t.Errorf("Get() = %v, want %v", store.Get("end"), "yes")
	}
	store.Close()
}
