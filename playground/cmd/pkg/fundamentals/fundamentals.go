package fundamentals

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
)

var Hello string = "Hello"

// Rune Sample One
func RuneSampleOne() {
	fmt.Println("\n---> Rune Sample One")
	s := "hello 👋"
	r := []rune(s)
	fmt.Println("rune --> ", r)

	for _, r := range "👋🌍" {
		fmt.Printf("%c = %U\n", r, r)
	}
}

// Defer Sample One
func DeferSampleOne() {
	fmt.Println("\n----> Mutex Examples")
	defer fmt.Println("A")
	fmt.Println("B")
}

// Defer Inside Loop
func DeferInsideLoop() {
	fmt.Println("\n---> Loop with Defer:")
	for i := 0; i < 3; i++ {
		defer fmt.Println("-> deferred:", i)
	}
}

// Arrays Sample One
func ArraysSampleOne() {
	fmt.Println("---> Array Sample One")

	var fruitList [4]string

	fruitList[0] = "Apple"
	fruitList[1] = "Tomato"
	fruitList[2] = "Peach"

	fmt.Println("Fruit list is: ", fruitList)

	var vegList = [5]string{"potato", "beans", "mushroom"}
	fmt.Println("vegList is: ", vegList)
}

// Slice Sample One
func SlicesSampleOne() {
	fmt.Println("\n---> Slice Sample One")

	var fruitList = []string{"Apple", "Tomato", "Peach"}
	fruitList = append(fruitList, "Mango", "Banana")
	fmt.Printf("Type of fruitList is %T\n", fruitList)
	fmt.Println("Before Sliced, Fruit List: ", fruitList)

	fruitList = fruitList[1:3]
	fmt.Println("After Sliced, Fruit List: ", fruitList)

	highScores := make([]int, 4)
	highScores[0] = 234
	highScores[1] = 235
	highScores[2] = 236
	highScores[3] = 456
	highScores = append(highScores, 555, 666, 777)
	fmt.Println("highScores --> ", highScores)
}

// Channel Sample One
func ChannelBasicSyntax() {
	fmt.Println("\n---> Channel Basic Syntax")

	ch := make(chan int)

	go func() {
		ch <- 42
	}()

	val := <-ch
	fmt.Println("value -> ", val)
}

func sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Variadic Functions
func VariadicFunctionSample() {
	fmt.Println(sum(1, 2, 3))
	fmt.Println(sum(10, 20, 3, 40))
}

func greet(names ...string) {
	for _, name := range names {
		fmt.Println("Hello", name)
	}
}

func UnpackingVariadicSample() {
	people := []string{"Alice", "Bob", "Charlie"}

	greet("John", "Doe")
	greet(people...)
}

func BitWiseRightShiftOperator() {
	x := 8      // binary: 1000
	y := x >> 1 // right shift by 1 = 0100 (4)
	z := x >> 2 // right shift by 2 = 0010 (2)

	fmt.Println(y) // 4
	fmt.Println(z) // 2
}

func StrConvAtoiSample() {
	input := "123"
	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Println("Conversion error:", err)
	} else {
		fmt.Println("Converted number:", num)
	}
}

func ArrayFilterSample() {
	nums := []int{1, 2, 3, 4, 5}
	var evens []int

	for _, n := range nums {
		if n%2 == 0 {
			evens = append(evens, n)
		}
	}

	fmt.Println("evens :>> ", evens)
}

func MarshalSample() {
	type User struct {
		Name string
		Age  int
	}

	user := User{Name: "Alice", Age: 30}
	jsonBytes, _ := json.Marshal(user)
	fmt.Println("jsonBytes --> ", jsonBytes)
	fmt.Println("jsonBytes string --> ", string(jsonBytes))

	jsonStr := string(jsonBytes)
	var unmarhsal_user User
	err := json.Unmarshal([]byte(jsonStr), &unmarhsal_user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("unmarhsal_user --> ", unmarhsal_user.Name)
}

func TypeAssertSampleOne() {
	var val interface{} = "hello"
	str := val.(string)
	fmt.Println(strings.ToUpper(str))
}

func TypeAssertDecodingDynamicData() {
	var data interface{}
	json.Unmarshal([]byte(`{"name": "thuta"}`), &data)

	m := data.(map[string]interface{})
	fmt.Println("data -> ", m["name"])
}

func TypeAssertAccessingValuesFromChannel() {
	ch := make(chan interface{})
	ch <- 42

	val := <-ch
	num := val.(int)
	fmt.Print(num)
}
