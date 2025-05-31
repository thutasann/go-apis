package main

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

const User_Bucket_Name = "users"

// Realtime Database
func main() {
	db, err := bbolt.Open(".db", 0666, nil)
	if err != nil {
		log.Fatal(err)
	}
	user := map[string]string{
		"name": "thutasann",
		"age":  "23",
	}

	db.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucket([]byte(User_Bucket_Name))
		if err != nil {
			return err
		}

		id := uuid.New()

		for k, v := range user {
			if err := bucket.Put([]byte(k), []byte(v)); err != nil {
				return err
			}
		}

		if err := bucket.Put([]byte("id"), []byte(id.String())); err != nil {
			return err
		}

		return nil
	})

	userRes := make(map[string]string)

	if err := db.View(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte(User_Bucket_Name))
		if bucket == nil {
			return fmt.Errorf("bucket name not found: %s", User_Bucket_Name)
		}

		bucket.ForEach(func(k, v []byte) error {
			userRes[string(k)] = string(v)
			return nil
		})

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	fmt.Println(userRes)
}
