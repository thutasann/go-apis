package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func Stream_Reader_Example() {
	r := strings.NewReader("Learning is fun")

	buf := make([]byte, 4)

	for {
		n, err := r.Read(buf)
		fmt.Println(string(buf[:n]), err)
		if err != nil {
			fmt.Println("breaking out...")
			break
		}
	}
}

func Stream_Reader_Example_Two() {
	r := strings.NewReader("some io.Reader stream to be read\n")
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}
}
