package fundamentals

import "fmt"

type Person struct {
	Name string
	Age  int
}

func StructPointerSampleOne() {
	fmt.Println("-----> Struct Pointer Sample One")
	p1 := Person{Name: "John", Age: 30}
	fmt.Println("p1 --> ", p1)

	p2 := Person{}
	fmt.Println("p2 --> ", p2)

	p3 := &Person{Name: "Alice", Age: 25}
	fmt.Println("p3 --> ", p3)
}

// Value Receiver
func (p Person) Greet() {
	fmt.Println("Hello My Name is ", p.Name)
}

// Pointer Receiver
func (p *Person) HaveBirthday() {
	p.Age++
}

func StructPointerSampleTwo() {
	fmt.Println("-----> Struct Pointer Sample Two")
	p := Person{Name: "John", Age: 30}
	ptr := &p
	fmt.Println("Person Pointer Name --> ", ptr.Name) // John
	fmt.Println("Person Pointer Pointer --> ", &ptr)  // John
}
