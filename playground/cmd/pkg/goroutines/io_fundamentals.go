package goroutines

import (
	"fmt"
	"io"
	"os"
	"strings"
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

func TeeReaderSample() {
	src := strings.NewReader("Hello World")
	tee := io.TeeReader(src, os.Stdout)
	io.ReadAll(tee)
}

func TeeeReaderSampleTwo() {
	original := strings.NewReader("this is secret data\n")

	// TeeReader will duplicate everything read from original to stdout
	tee := io.TeeReader(original, os.Stdout)

	result, err := io.ReadAll(tee)
	if err != nil {
		panic(err)
	}

	fmt.Printf("\nCaptured result : %s", string(result))
}

func MultiReaderSample() {
	r1 := strings.NewReader("Part1\n")
	r2 := strings.NewReader("Part2\n")
	r3 := strings.NewReader("Part3\n")

	multi := io.MultiReader(r1, r2, r3)
	io.Copy(os.Stdout, multi)
}

func MultiWriterSample() {
	f, _ := os.Create("log.txt")
	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f)
	mw.Write([]byte("Logging this to stdout and file\n"))
}

func RedirectStdoutToAFile() {
	f, _ := os.Create("stdout_log.txt")
	defer f.Close()
	os.Stdout = f // redirect all stdout to the file
	println("This goes to stdout_log.txt instead of the terminal.")
}

// - uint8 is an unsigned 8-bit integer
// - range: 0 to 255
// - its literally 1 byte of memory (8 bits), no sign bit
func Uint8Sample() {
	var a uint8 = 65
	fmt.Println(a)
	fmt.Printf("%c\n", a)
}

// byte -> a slice of uint8 values
func ByteSample() {
	var b []byte = []byte{72, 101, 108, 108, 111}
	fmt.Println(string(b)) // Hello
	fmt.Println(b)         // [72 101 108 108 111]
}
