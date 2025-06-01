package hopper

import (
	"fmt"

	"github.com/google/uuid"
	"go.etcd.io/bbolt"
)

const (
	defaultDBName = "default"
)

// Data Map
type M map[string]string

// Hooper struct
type Hopper struct {
	db *bbolt.DB // Hopper Database
}

// Collection Struct
type Collection struct {
	*bbolt.Bucket // Collection's Bucket
}

// Initialize New Hopper
func New() (*Hopper, error) {
	dbName := fmt.Sprintf("%s.hopper", defaultDBName)

	db, err := bbolt.Open(dbName, 0666, nil)
	if err != nil {
		return nil, err
	}
	return &Hopper{
		db: db,
	}, nil
}

// Create New collection
func (h *Hopper) CreateCollection(name string) (*Collection, error) {
	tx, err := h.db.Begin(true)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	bucket := tx.Bucket([]byte(name))
	if bucket != nil {
		return &Collection{Bucket: bucket}, nil
	}

	bucket, err = tx.CreateBucket([]byte(name))
	if err != nil {
		return nil, err
	}

	return &Collection{Bucket: bucket}, nil
}

// Insert Data
func (h *Hopper) Insert(collName string, data M) (uuid.UUID, error) {
	id := uuid.New()

	coll, err := h.CreateCollection(collName)
	if err != nil {
		return id, err
	}

	h.db.Update(func(tx *bbolt.Tx) error {
		for k, v := range data {
			if err := coll.Put([]byte(k), []byte(v)); err != nil {
				return err
			}
		}

		if err := coll.Put([]byte("id"), []byte(id.String())); err != nil {
			return err
		}

		return nil
	})

	return id, nil
}

// Select Data
func (h *Hopper) Select(coll string, k string, query interface{}) {

}
