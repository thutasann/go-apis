package fundamentals

import (
	"fmt"
)

var Hello string = "Hello"

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
