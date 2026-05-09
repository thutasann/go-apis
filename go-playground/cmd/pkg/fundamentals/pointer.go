package fundamentals

import (
	"encoding/json"
	"fmt"
)

func PointerSampleOne() {
	fmt.Println("---> Pointer Sample One ")
	num := 10
	ptr := &num
	fmt.Println("num --> ", num)    // 10
	fmt.Println("&num --> ", &num)  // 0x14000112008
	println("ptr --> ", ptr)        // 0x14000112008
	println("ptr value --> ", *ptr) // 10
}

func modifyFn(x *int) {
	*x = 100
}

func ModifyPointerFunctionSample() {
	println("----> Modify Pointer Function")
	num := 10
	println("Before : ", num)

	modifyFn(&num)
	println("After : ", num)
}

type User struct {
	Name string
	Age  int
}

// Pointer receiver modifies the struct
func (u *User) birthday() {
	u.Age++
}

// Pointers to Structs & Methods
func PointerStructSampleOne() {
	user := User{Name: "Alice", Age: 25}
	user.birthday()
	fmt.Println(user)
}

func double_pointer_modify(ptr **int) {
	**ptr = 50
}

// Modify the original value with Double Pointer
func DoublePointer() {
	println("------> Double Pointer")
	num := 10
	ptr := &num
	ptr2 := &ptr // pointer to a pointer

	double_pointer_modify(ptr2)
	fmt.Println(num)
}

func arr_slice_modify(arr *[3]int) {
	(*arr)[0] = 99
}

func ArraySliceModify() {
	nums := [3]int{1, 2, 3}
	arr_slice_modify(&nums)
	fmt.Println(nums)
}

func StringVsPointerString() {
	fmt.Println("---> String vs Pointer String")
	var normalString string
	fmt.Println("normalString --> ", normalString) // ""

	var s *string
	fmt.Println(s) // nil

	value := "hello"
	s = &value
	fmt.Println(*s) // "hello"
}

// Modify Without Pointer Sample
func WithoutPointerSample() {
	fmt.Println("\n---> WithoutPointerSample")
	setToZero := func(n int) {
		n = 0
	}
	x := 5
	setToZero(x)
	fmt.Println(x) // 5 - not changed
}

// Modify With Pointer Sample
// - x is a house.
// - &x is the address of the house.
// - *p is you going to that address and changing what's inside the house.
func WithPointerSample() {
	fmt.Println("\n---> With PointerSample")

	setToZero := func(n *int) {
		*n = 0
	}
	x := 5
	setToZero(&x)
	fmt.Println(x) // 0 - changed!
}

type Config struct {
	Port  int
	Debug bool
}

// Modifying Config
func ModifyingConfig() {
	fmt.Println("\n---> Modifying Config")
	var cfg Config
	loadConfig(&cfg)
	fmt.Println(cfg)
}

func loadConfig(cfg *Config) {
	cfg.Port = 8000
	cfg.Debug = true
}

// In Place Sorting
func InPlaceSorting() {
	fmt.Println("\n---> In-Place Sorting")

	swap := func(a, b *int) {
		*a, *b = *b, *a
	}
	x, y := 5, 10
	swap(&x, &y)
	fmt.Println(x, y)
}

type Food struct {
	Name *string `json:"name"`
}

func StringVsPointerStringTwo() {
	// JSON where `name` is completely missing
	jsonMissing := `{}`

	// JSON where `name` is explicitly empty
	jsonEmtpy := `{"name":""}`

	var food1, food2 Food

	_ = json.Unmarshal([]byte(jsonMissing), &food1)
	_ = json.Unmarshal([]byte(jsonEmtpy), &food2)

	// food1.name is nil because 'name' was not present
	if food1.Name == nil {
		fmt.Println("food1: name was not sent")
	} else {
		fmt.Println("food1: name=", *food1.Name)
	}

	// food2.Name is a pointer to empty string
	if food2.Name == nil {
		fmt.Println("food2: name was not sent")
	} else {
		fmt.Printf("food2: name = '%s'\n", *food2.Name)
	}
}

type Food2 struct {
	Name string
}

// INCORRECT: passing by value
func fill(f Food2) {
	f.Name = "Pizza" // modify only the copy
}

// CORRECT: passing by pointer
func fillPtr(f *Food2) {
	f.Name = "Pizza" // modifies the actual original
}

// Comparison Between to pass by value (copy) vs pass by pointer
func PassPointerSample() {
	var food = Food2{Name: "Burger"}
	fill(food)
	fmt.Println("Food fill --> ", food)

	fillPtr(&food)
	fmt.Println("Food fillPtr --> ", food)
}

type PUser struct {
	Name string
}

func changeName(u PUser) {
	u.Name = "Bob" // chagnes only to the copy
	fmt.Println("u (changeName) --> ", u)
}

// When to Pass value (copy)
//
// - You don’t want the original to be modified.
//
// - The data is small and cheap to copy.
//
// - You want to isolate state changes inside the function.
func WhenToPassValue() {
	original := PUser{Name: "Alice"}
	changeName(original)
	fmt.Println("original --> ", original) // still "Alice"
}
