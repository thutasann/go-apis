package demos

import "fmt"

func DeferDemo() {
	fmt.Println("\n---> Defer Demo")
	fmt.Println("Return test:", testReturn())
	fmt.Println("Named return test:", testNamedReturn())
	fmt.Println("Panic test:")
	testPanic()
}

func testReturn() int {
	defer fmt.Println("-> defer: testReturn")
	return 42
}

func testNamedReturn() (res int) {
	defer func() {
		fmt.Println("-> defer: testNamedRetrun (before res change)", res)
		res = 100
		fmt.Println("-> defer: testNamedReturn (after res change)", res)
	}()
	res = 10
	return
}

func testPanic() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("-> recovered from panic:", r)
		}
	}()
	panic("🔥 something went wrong!")
}
