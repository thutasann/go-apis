package goroutines

import (
	"fmt"
	"io"
	"os"
)

func ReadingFileSample() {
	file, err := os.Open("sample.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	buf := make([]byte, 100)
	for {
		n, err := file.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}
		fmt.Println(string(buf[:n]))
	}
}

func WritingFileSample() {
	file, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	message := "this is a log entry\n"
	n, err := file.Write([]byte(message))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Wrote %d bytes\n", n)
}

func CopyFromReadertoWriter() {
	src, _ := os.Open("sample.txt")
	defer src.Close()

	dst, _ := os.Create("copy.txt")
	defer dst.Close()

	io.Copy(dst, src)
}
