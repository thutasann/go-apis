package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// Database struct
type Database struct {
	db *sql.DB
}

// Initialize a new Database
func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "postgresql://root:password@localhost:5433/go-chat?sslmode=disable")

	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

// Close the database
func (d *Database) Close() {
	d.db.Close()
}

// Get the database
func (d *Database) GetDB() *sql.DB {
	return d.db
}
