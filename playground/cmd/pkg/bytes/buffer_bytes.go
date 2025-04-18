package bytes

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

// Buffer Sample One
func BufferSampleOne() {
	fmt.Println("----> Buffer Sample One")
	var b bytes.Buffer
	b.WriteString("Hello, ")
	b.WriteString("World!")
	fmt.Println(b.String())
}

// Bytes Samples
func BytesSamples() {
	fmt.Println("----> Bytes Samples")
	buf := new(bytes.Buffer)
	buf.Write([]byte("Foo"))
	buf.WriteString("Bar")
	fmt.Println(buf.Len())
	fmt.Println(buf.String())
}

// Buffered Network I/O
//
// Use case : Optimize TCP read/write to avoid syscall overhead
func BufferedNetworkIO() {
	ln, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept: %v", err)
			continue
		}
		go handleConn(conn)
	}
}

// Handle Connection
func handleConn(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		fmt.Printf("Received: %s", line)

		writer.WriteString("Echo: " + line)
		writer.Flush()
	}
}

type Event struct {
	UserID string `json:"user_id"`
	Action string `json:"action"`
}

// Streaming JSON to a Buffer (Logging)
func StreamingJSONToBuffer() {
	fmt.Println("----> Streaming JSON to a Buffer ")
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)

	events := []Event{
		{"123", "login"},
		{"123", "click"},
		{"123", "logout"},
	}

	for _, e := range events {
		encoder.Encode(e)
	}

	fmt.Println(buf.String())
}

type BytesClosingBuffer struct {
	*bytes.Buffer
	io.Closer
}

func NewBytesClosingBuffer() *BytesClosingBuffer {
	return &BytesClosingBuffer{
		Buffer: new(bytes.Buffer),
	}
}

func (b *BytesClosingBuffer) Close() error {
	fmt.Println("closing...")
	return nil
}

// Write To Sample
func IOWriteToSample() {
	fmt.Println("----> IO Write To Sample")
	buf := NewBytesClosingBuffer()

	if err := writeTo(buf, []byte("hello world")); err != nil {
		log.Fatal(err)
	}
	fmt.Println(buf)
}

func writeTo(wc io.WriteCloser, msg []byte) error {
	defer wc.Close()
	_, err := wc.Write(msg)
	return err
}

// Buffered File Writer
func BufferedFileWriter() {
	file, _ := os.Create("output.txt")
	defer file.Close()

	writer := bufio.NewWriter(file)

	for i := 0; i < 100000; i++ {
		writer.WriteString("line\n")
	}

	writer.Flush()
}
