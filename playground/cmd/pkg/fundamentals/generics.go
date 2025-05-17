package fundamentals

import (
	"fmt"
)

// ----------- 🧩 Beginner: Generic Function

// Generic Print Function
// T is a type parameter - can be any type
func print[T any](value T) {
	fmt.Println(value)
}

// Generic Sample One
func GenericSampleOne() {
	print("Hello")
	print(123)
	print(true)
}

// ----------- 🧮 Generic Function with Constraints

type Number interface {
	int | int64 | float64
}

func add[T Number](a, b T) T {
	return a + b
}

// Generic Function with Constraints
func GenericWithConstraints() {
	fmt.Println(add(1, 2))
	fmt.Println(add(1.5, 1.2))
}

// ----------- 📦 Generic Struct

type Box[T any] struct {
	Value T
}

type GenericUser struct {
	Name string
}

func GenericStructSample() {
	intBox := Box[int]{Value: 12}
	stringBox := Box[string]{Value: "String"}
	booleanBox := Box[bool]{Value: true}
	userBox := Box[GenericUser]{Value: GenericUser{Name: "Thuta"}}

	fmt.Println(intBox.Value)
	fmt.Println(stringBox.Value)
	fmt.Println(booleanBox.Value)
	fmt.Println(userBox.Value.Name)
}

// ----------- 📚  Generic Slice Filter Function

func GenericFilter[T any](data []T, test func(T) bool) []T {
	var result []T
	for _, v := range data {
		if test(v) {
			result = append(result, v)
		}
	}

	return result
}

func GenericFilterUasge() {
	nums := []int{1, 2, 3, 4, 5}
	even := GenericFilter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println(even)
}

// ----------- 🧠  Advanced: Custom Constraints

type Comparable interface {
	~int | ~string // tilde means "any type with underlying int/string"
}

func Max[T Comparable](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// ----------- 🔬 Generic Map Function

func Map[T any, R any](data []T, mapper func(T) R) []R {
	var result []R

	for _, v := range data {
		result = append(result, mapper(v))
	}

	return result
}

func GenericMap() {
	names := []string{"Go", "Rust"}
	lengths := Map(names, func(s string) int {
		return len(s)
	})
	fmt.Println(lengths) // [2, 4]
}

// ----------- 🧪  Real-world: API Response Wrapper
type ApiResponse[T any] struct {
	Success bool `json:"success"`
	Data    T    `json:"data"`
}

func WrapSuccess[T any](data T) ApiResponse[T] {
	return ApiResponse[T]{Data: data, Success: true}
}
