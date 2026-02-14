package main

import (
	"errors"
	"fmt"
)

func process_int(i int) error {
	if i%2 == 0 {
		return errors.New("only odd numbers allowed")
	}
	return nil
}

func check_error(e error) {
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println("operation successful")
}

func Errors_One() {
	err := process_int(3)
	check_error(err)
}
