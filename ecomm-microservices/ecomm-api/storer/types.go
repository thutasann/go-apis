package storer

import "time"

type Product struct {
	ID           int64      `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Image        string     `json:"image" db:"image"`
	Category     string     `json:"category" db:"category"`
	Description  string     `json:"description" db:"description"`
	Rating       int64      `json:"rating" db:"rating"`
	NumReviews   int64      `json:"num_reviews" db:"num_reviews"`
	Price        float64    `json:"price" db:"price"`
	CountInStock int64      `json:"count_in_stock" db:"count_in_stock"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at" db:"updated_at"`
}

type Order struct {
	ID            int64      `db:"id"`
	PaymentMethod string     `db:"payment_method"`
	TaxPrice      float64    `db:"tax_price"`
	ShippingPrice float64    `db:"shipping_price"`
	TotalPrice    float64    `db:"total_price"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     *time.Time `db:"updated_at"`
	Items         []OrderItem
}

type OrderItem struct {
	ID        int64   `db:"id"`
	Name      string  `db:"name"`
	Quantity  int64   `db:"quantity"`
	Image     string  `db:"image"`
	Price     float64 `db:"price"`
	ProductID int64   `db:"product_id"`
	OrderID   int64   `db:"order_id"`
}
