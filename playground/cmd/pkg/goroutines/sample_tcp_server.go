package goroutines

import (
	"fmt"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024) // create buffer to read data
	n, err := conn.Read(buffer)  // read data into buffer
	if err != nil {
		fmt.Println("Error reading: ", err)
	}

	message := string(buffer[:n])
	fmt.Println("Received: ", message)

	conn.Write([]byte("Hello back!\n")) // send reply back
}

func SimpleTCPServer() {
	listener, err := net.Listen("tcp", ":9000") // listen on tcp port 9000
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("🚀 Server listening on port 9000...")

	for {
		conn, err := listener.Accept() // Waits (blocking) until a client connects
		if err != nil {
			fmt.Println("Error accepting:", err)
			continue
		}

		go handleConnection(conn)
	}
}
