package caskdb

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestDiskStore_Get(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")
	store.Set("name", "jojo")
	if val, _ := store.Get("name"); val != "jojo" {
		t.Errorf("Get() = %v, want %v", val, "jojo")
	}
}

func TestDiskStore_GetInvalid(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")
	if val, _ := store.Get("some key"); val != "" {
		t.Errorf("Get() = %v, want %v", val, "")
	}
}

func TestDiskStore_SetWithPersistence(t *testing.T) {
	store, err := NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	defer os.Remove("test.db")

	tests := map[string]interface{}{
		"crime and punishment": "dostoevsky",
		"anna karenina":        "tolstoy",
		"hamlet":               "shakespeare",
		"floatval":             2121.2123,
		"boolval":              true,
		"runeval":              UglyRune('&'),
		"intval":               -12.212,
	}
	for key, val := range tests {
		store.Set(key, val)
		actualVal, _ := store.Get(key)
		if actualVal != val {
			t.Errorf("Get() = %v, want %v", actualVal, val)
		}
	}
	type Person struct {
		Name string
		Age  int
	}
	p := &Person{
		Name: "Sebastian Vettel",
		Age:  36,
	}
	bytes, _ := json.Marshal(p)
	store.Set("mystruct", bytes)
	store.Close()

	store, err = NewDiskStore("test.db")
	if err != nil {
		t.Fatalf("failed to create disk store: %v", err)
	}
	for key, val := range tests {
		actualVal, _ := store.Get(key)
		if actualVal != val {
			t.Errorf("Get() = %v, want %v", actualVal, val)
		}
	}

	store.SetX("expire", -120.212, 2*time.Second)

	value, _ := store.Get("mystruct")
	mystruct := &Person{}
	err = json.Unmarshal(value.([]byte), mystruct)
	if err != nil {
		t.Fatalf("couldn't unmarshal struct types: %v", err)
	}
	if mystruct.Age != p.Age && mystruct.Name != p.Name {
		t.Errorf("Get() = %v, want %v", value.(*Person), p)
	}

	//check for expired key
	time.Sleep(3 * time.Second)
	_, err = store.Get("expire")
	if err != nil {
		if err != ErrKeyNotFound {
			t.Error("expected key to be expired")
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
		actualVal, _ := store.Get(dkeys)
		if actualVal != TombStoneVal {
			t.Errorf("Get() = %s, want %s", actualVal, TombStoneVal)
		}
	}
	store.Close()
}
