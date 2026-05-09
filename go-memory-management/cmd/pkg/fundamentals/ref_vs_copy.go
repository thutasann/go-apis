package fundamentals

import "fmt"

type CopyValueUser struct {
	Name string
	Age  int
}

// Struct Copy (Value)
func StructCopySample() {
	p1 := CopyValueUser{Name: "John", Age: 30}
	p2 := p1 // copied by value

	p2.Name = "Jane"

	fmt.Println(p1.Name) // John
	fmt.Println(p2.Name) // Jane
}

// Struct Reference (Pointer)
func StructReferenceSample() {
	p1 := CopyValueUser{Name: "John", Age: 30}
	p2 := &p1 // pointer

	p2.Name = "Jane"

	fmt.Println(p1.Name)
	fmt.Println(p2.Name)
}

// Because slice header is copied, but it still points to the same backing array.
func SliceCopyPitfall() {
	data := []int{1, 2, 3}
	copyData := data
	copyData[0] = 99
	fmt.Println(data[0])
}

// Now They are truly separate
func ProperSliceCopy() {
	data := []int{1, 2, 3}
	copyData := make([]int, len(data))
	copy(copyData, data)

	copyData[0] = 99
	fmt.Println(data[0])
}

// ----------  Function Argument Pitfall

type Config struct {
	Timeout int
}

func pitfall_mutate(cfg Config) {
	cfg.Timeout = 100
}

// proper: pass by pointer
func proper_mutate(cfg *Config) {
	cfg.Timeout = 100
}

func PitfallMutateSample() {
	config := Config{Timeout: 10}
	pitfall_mutate(config)
	fmt.Println(config.Timeout) // 10 - not mutate
}

func ProperMutateSample() {
	config := Config{Timeout: 10}
	proper_mutate(&config)
	fmt.Println(config.Timeout) // 100 - mutate
}

// ----------  Reference Impact: Performance and Garbage

type Point struct {
	X, Y int
}

// P is copied
// Fine for small structs (no heap, no GC pressure)
func DoSomething(p Point) {
	fmt.Println(p.X)
}

type Large struct {
	Data [1 << 20]byte // 1MB
}

// Only pointer copied
// Huge win: avoids copying large memory on every call
func DoSomethingLarge(p *Large) {
	fmt.Println(p.Data)
}

// ---------- Visualization of Reference vs Copy

type RefVsCopyUser struct {
	Name string
}

func RefVsCopyVisualization() {
	// copy
	a := RefVsCopyUser{Name: "Alice"}
	b := a // indepedent
	b.Name = "Bob"

	// Reference
	c := &a
	c.Name = "Charlie"

	fmt.Println(a.Name) // Charlie (because `c` points to `a`)
	fmt.Println(b.Name) // Bob (Indepedent copy)
}
