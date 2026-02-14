package main

import (
	"fmt"
	"os"
)

func Read_File_Sample() {
	path := FILE_PATH
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(data))
}

func OS_Open_Sample() {
	path := FILE_PATH
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Open Error: ", err)
	}

	b := make([]byte, 4)
	for {
		n, err := file.Read(b)
		if err != nil {
			fmt.Println("Error: ", err)
			break
		}
		fmt.Println(string(b[0:n]))
	}
}
