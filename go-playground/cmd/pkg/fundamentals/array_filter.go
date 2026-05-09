package fundamentals

import "fmt"

// Filter the Arrray Utility
func Filter[T any](items []T, test func(T) bool) []T {
	var filtered []T

	for _, item := range items {
		if test(item) {
			filtered = append(filtered, item)
		}
	}

	return filtered
}

// Filter Map Utility
func FilterMap[T any, R any](items []T, f func(T) R) []R {
	var result []R
	for _, item := range items {
		result = append(result, f(item))
	}
	return result
}

func ArrayFilterGenericSample() {
	nums := []int{1, 2, 3, 4, 5}
	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Println("evens :>> ", evens)

	words := []string{"go", "node", "rust"}
	gos := Filter(words, func(s string) bool { return s == "go" })
	fmt.Println("gos :>> ", gos) // [go]
}

func FilteringCustomStruct() {
	type ArrayUser struct {
		Name string
		Age  int
	}

	users := []ArrayUser{
		{"Alice", 30},
		{"Bob", 17},
		{"Charlie", 25},
	}

	adults := Filter(users, func(u ArrayUser) bool {
		return u.Age >= 18
	})

	for _, a := range adults {
		fmt.Println("adult :>> ", a.Name)
	}

	names := FilterMap(adults, func(u ArrayUser) string {
		return u.Name
	})
	fmt.Println(names) // ["Alice", "Charlie"]
}
