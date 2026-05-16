package main

import (
	"fmt"
	"os"
)

func Append_File_One() {
	file, err := os.OpenFile(FILE_PATH, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer file.Close()

	_, err = file.WriteString("Hope you had a good day!")
	if err != nil {
		fmt.Println(err)
	}
}
