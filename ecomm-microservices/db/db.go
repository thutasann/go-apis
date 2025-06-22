package db

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Database struct
type Database struct {
	db *sqlx.DB
}

// Initialize a New Database
func NewDatabase() (*Database, error) {
	db, err := sqlx.Open("mysql", "root:password@tcp(localhost:3306)/ecom?parseTime=true")
	if err != nil {
		return nil, fmt.Errorf("error opening database: %v", err)
	}

	return &Database{db: db}, nil
}

// Close Database
func (d *Database) Close() error {
	return d.db.Close()
}

// Get Database
func (d *Database) GetDB() *sqlx.DB {
	return d.db
}
