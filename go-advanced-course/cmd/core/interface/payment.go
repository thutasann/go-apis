package main

import (
	"fmt"
	"time"
)

type Payment interface {
	Process(amount float64) error
	Refund(amount float64) error
}

type Stripe struct {
	APIKey string
}

func (s Stripe) Process(amount float64) error {
	fmt.Printf("[Stripe] Processing payment of $%.2f using API Key: %s\n", amount, s.APIKey)
	time.Sleep(500 * time.Millisecond)
	fmt.Println("[Stripe] Payment successful!")
	return nil
}

func (s Stripe) Refund(amount float64) error {
	fmt.Printf("[Stripe] Refunding $%.2f using API Key: %s\n", amount, s.APIKey)
	time.Sleep(300 * time.Millisecond)
	fmt.Println("[Stripe] Refund successful!")
	return nil
}

type PayPal struct {
	ClientID string
}

func (p PayPal) Process(amount float64) error {
	fmt.Printf("[PayPal] Processing payment of $%.2f with ClientID: %s\n", amount, p.ClientID)
	time.Sleep(400 * time.Millisecond)
	fmt.Println("[PayPal] Payment successful!")
	return nil
}

func (p PayPal) Refund(amount float64) error {
	fmt.Printf("[PayPal] Refunding $%.2f with ClientID: %s\n", amount, p.ClientID)
	time.Sleep(200 * time.Millisecond)
	fmt.Println("[PayPal] Refund successful!")
	return nil
}

func ExecutePayment(p Payment, amount float64) {
	fmt.Println("Starting payment execution...")
	err := p.Process(amount)
	if err != nil {
		fmt.Println("Payment failed:", err)
		return
	}
	fmt.Println("Payment executed successfully!")
}

func ExecuteRefund(p Payment, amount float64) {
	fmt.Println("Starting refund execution...")
	err := p.Refund(amount)
	if err != nil {
		fmt.Println("Refund failed:", err)
		return
	}
	fmt.Println("Refund executed successfully!")
}

func Payment_Interface_Sample() {
	stripe := Stripe{APIKey: "sk_test_12345"}
	paypal := PayPal{ClientID: "client_abc"}

	ExecutePayment(stripe, 99.99)
	fmt.Println("------")
	ExecuteRefund(stripe, 50.00)
	fmt.Println("------")
	ExecutePayment(paypal, 120.50)
	fmt.Println("------")
	ExecuteRefund(paypal, 30.00)
}
