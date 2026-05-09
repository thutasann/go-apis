package goroutines

import (
	"fmt"
	"sort"
	"sync"
)

func MapSampleOne() {
	ages := map[string]int{
		"Alice": 30,
		"Bob":   25,
	}
	fmt.Println(ages)
	fmt.Println("Alice age", ages["Alice"])
}

func MapDeclaringAndUsing() {
	m := make(map[string]int)

	m["apple"] = 4

	fmt.Println(m["apple"])

	delete(m, "apple")

	val, ok := m["apple"]
	if ok {
		fmt.Println("Exists:", val)
	} else {
		fmt.Println("not found")
	}
}

func MapOfStructs() {
	type User struct {
		Name string
		Age  int
	}

	users := map[string]User{
		"user1": {Name: "Alice", Age: 22},
		"user2": {Name: "Bob", Age: 30},
	}

	fmt.Println(users)
}

func MapOfSlices() {
	m := map[string]map[string]int{
		"john": {"math": 90, "sience": 85},
	}
	fmt.Println(m)
}

func SortingAMap() {
	m := map[string]int{"b": 2, "a": 1, "c": 3}
	var keys []string

	for k := range m {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, m[k])
	}
}

func MapWithInterface() {
	m := map[string]interface{}{
		"name": "Alice",
		"age":  30,
		"meta": map[string]string{"city": "NY"},
	}
	fmt.Println(m)
}

func ConcurrencySafeMap() {
	var sm sync.Map
	sm.Store("a", 1)
	value, ok := sm.Load("a")
	fmt.Println(value, ok)
}

func MapWithLock() {
	var mu sync.RWMutex
	m := make(map[string]int)

	// safe write
	mu.Lock()
	m["a"] = 1
	mu.Unlock()

	// safe read
	mu.RLock()
	val := m["a"]
	mu.RUnlock()
	fmt.Println(val)
}

func CountingFrequency() {
	count := make(map[string]int)
	words := []string{"go", "go", "lang"}
	for _, w := range words {
		count[w]++
	}
	fmt.Println("count -->", count)
}

func BasicMapLoop() {
	m := map[string]int{
		"Alice": 25,
		"Bob":   30,
	}

	for key, value := range m {
		fmt.Printf("%s is %d years old\n", key, value)
	}
}
