package paymentinterface

import (
	"fmt"
	"log"
)

// PaymentProvider defines the behavior any payment provider must implement
type PaymentProvider interface {
	Pay(amount float64) (string, error)
	Refund(transactionID string) error
}

// Stripe implements PaymentProvider
type Stripe struct {
	APIKey string
}

func (s *Stripe) Pay(amount float64) (string, error) {
	fmt.Printf("Stripe charged $%.2f\n", amount)
	return "strip_txn_123", nil
}

func (s *Stripe) Refund(transactionID string) error {
	fmt.Printf("Stripe refunded %s\n", transactionID)
	return nil
}

// PayPal implements PaymentProvider
type PayPal struct {
	ClientID string
}

func (p *PayPal) Pay(amount float64) (string, error) {
	fmt.Printf("PayPal charged $%.2f\n", amount)
	return "paypal_txn_456", nil
}

func (p *PayPal) Refund(transactionID string) error {
	fmt.Printf("PayPal refunded %s\n", transactionID)
	return nil
}

// PaymentService uses any provider that satisfies PaymentProvider
type PaymentService struct {
	Provider PaymentProvider
}

func (ps *PaymentService) Checkout(amount float64) {
	txID, err := ps.Provider.Pay(amount)
	if err != nil {
		log.Println("Payment failed:", err)
		return
	}
	fmt.Println("Payment successful, txID:", txID)
}

func PaymentInterfaceExample() {
	stripe := &Stripe{APIKey: "sk_test_123"}
	paypal := &PayPal{ClientID: "client_xyz"}

	service := PaymentService{Provider: stripe}
	service.Checkout(99.99)

	service.Provider = paypal
	service.Checkout(49.50)
}
