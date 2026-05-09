package fundamentals

import "fmt"

// x in stackMemory() is cleaned up after the function ends.
func stackMemory() int {
	x := 42 // stored on stack (short-lived)
	return x
}

// x in heapMemory() survives because it's returned as a pointer.
func heapMemory() *int {
	x := 42
	return &x
}

// Stack vs Heap
func StackVsHeap() {
	int_x := stackMemory()
	pt_x := heapMemory()
	fmt.Println(int_x)
	fmt.Println(pt_x)
}

// ------------- Caching with Memory Management

type CacheUser struct {
	ID   int
	Name string
}

var userCache = make(map[int]*CacheUser)

// &User{...} ensures the object stays on the heap
// Caching avoids repeated allocations
// need to be aware of memory leaks if you never clear the cache
func getUser(id int) *CacheUser {
	if user, ok := userCache[id]; ok {
		return user
	}

	user := &CacheUser{ID: id, Name: "John Doe"} // Allocated on Heap
	userCache[id] = user
	return user
}

func CachingWithMemoryManagement() {
	user := getUser(1)
	fmt.Println(user.Name)
}
