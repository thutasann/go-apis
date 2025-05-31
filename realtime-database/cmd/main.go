package main

import (
	"fmt"
	"log"

	"github.com/thutasann/godb/hopper"
)

// Realtime Database
func main() {
	db, err := hopper.New()
	if err != nil {
		log.Fatal(err)
	}

	user := map[string]string{
		"name": "thutasann",
		"age":  "23",
	}

	id, err := db.Insert("user", user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("id :>>", id)
}
