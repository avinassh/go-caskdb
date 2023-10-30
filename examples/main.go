package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	caskDB "github.com/avinassh/go-caskdb"
)

func main() {
	store, err := caskDB.NewDiskStore("test.db")
	if err != nil {
		log.Fatalf("failed to create disk store: %v", err)
		os.Exit(-1)
	}
	defer store.Close()
	defer os.Remove("test.db")

	store.Set("screwderia", "charles leclrec") // cause ferrari screws everyone
	store.Set("redbull", "max verstappen")
	store.Set("mercedes", "lewis hamilton")
	store.Set("mclaren", "lando norris")

	//yeet lando
	store.Delete("mclaren")
	val, err := store.Get("mclaren")
	if val == "" && errors.Is(err, caskDB.ErrKeyNotFound) {
		fmt.Println("mclaren dropped lando norris for good!")
	}

	rbDriver, _ := store.Get("redbull")
	fmt.Printf("%s drives for redbull racing!", rbDriver)
}
