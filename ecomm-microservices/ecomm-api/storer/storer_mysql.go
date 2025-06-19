package storer

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type MySQLStorer struct {
	db *sqlx.DB
}

// Initialize a new MySQL Storer
func NewMySQLStorer(db *sqlx.DB) *MySQLStorer {
	return &MySQLStorer{db: db}
}

func (ms *MySQLStorer) CreateProduct(ctx context.Context, p *Product) {

}
