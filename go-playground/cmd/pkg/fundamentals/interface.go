package fundamentals

import (
	"fmt"
)

type Speaker interface {
	Speak() string
}

type Dog struct{}

func (d Dog) Speak() string {
	return "Woof!"
}

type Cat struct{}

func (c Cat) Speak() string {
	return "Meow!"
}

func SaySomething(s Speaker) {
	fmt.Println(s.Speak())
}

func PrintAnything(i interface{}) {
	fmt.Println(i)
}

// Interface Struct Usage Sample 1
func InterfaceStructUsage() {
	fmt.Println("----> Interface and Struct Usage")
	SaySomething(Dog{}) // Woof!
	SaySomething(Cat{}) // Meow!

	PrintAnything(42)
	PrintAnything("hello")
	PrintAnything(true)
}

// ----- Multiple Method
type Shape interface {
	Area() float64
	Perimeter() float64
}

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return 3.14 * c.Radius * c.Radius
}

func (c Circle) Perimeter() float64 {
	return 2 * 3.14 * c.Radius
}

// ----- Custom Logger

type Logger interface {
	Log(msg string)
}

type ConsoleLogger struct{}

func (c ConsoleLogger) Log(msg string) {
	fmt.Println("[console]", msg)
}

func DoSomething(logger Logger) {
	logger.Log("doing work")
}

// -------------- PaymentGateway --------------

// Payment Gateway Interface
type PaymentGateway interface {
	Charge(amount float64) error
}

type Stripe struct{}

func (s Stripe) Charge(amount float64) error {
	fmt.Printf("Charging $%.2f using Stripe...\n", amount)
	return nil
}

type Paypal struct{}

func (p Paypal) Charge(amount float64) error {
	fmt.Printf("Charging $%.2f using PayPal...\n", amount)
	return nil
}

func ProcessPayment(gateway PaymentGateway, amount float64) {
	err := gateway.Charge(amount)
	if err != nil {
		fmt.Println("Payment failed!", err)
	}
	fmt.Println("Payment successful!")
}

func InterfacePaymentGatewaySample() {
	fmt.Println("\n----> Interface Payment Gateway Sample")

	var gateway PaymentGateway

	gateway = Stripe{}
	ProcessPayment(gateway, 49.99)

	gateway = Paypal{}
	ProcessPayment(gateway, 29.99)
}
