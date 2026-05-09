package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

var store = make(map[string]string)
var mu sync.RWMutex

// Redis-like server
// - bufio.Reader to buffer incoming data.
// - []byte with bytes.Buffer behavior.
// - sync.RWMutex to safely read/write shared memory.
func main() {
	ln, err := net.Listen("tcp", ":6380")
	if err != nil {
		panic(err)
	}
	fmt.Println("Redis-like server running on :6380")
	for {
		conn, _ := ln.Accept()
		go handleConn(conn)
	}
}

// Handle Connection
func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		req, err := parseRESP(reader)
		if err != nil {
			if err != io.EOF {
				fmt.Println("ERR:", err)
			}
		}
		resp := handleCommand(req)
		conn.Write([]byte(resp))
	}
}

// Parse RESP
func parseRESP(r *bufio.Reader) ([]string, error) {
	line, err := r.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if line[0] != '*' {
		return nil, fmt.Errorf("invalid array")
	}

	count, _ := strconv.Atoi(strings.TrimSpace(line[1:]))
	result := make([]string, 0, count)

	for i := 0; i < count; i++ {
		_, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}

		data, err := r.ReadString('\n')
		if err != nil {
			return nil, err
		}
		result = append(result, strings.TrimSpace(data))
	}
	return result, nil
}

// Handle Command
func handleCommand(args []string) string {
	if len(args) == 0 {
		return "-ERR empty command\r\n"
	}

	cmd := strings.ToUpper(args[0])

	switch cmd {
	case "PING":
		return "+PONG\r\n"

	case "SET":
		if len(args) != 3 {
			return "-ERR wrong number of arguments for SET\r\n"
		}
		key, value := args[1], args[2]
		mu.Lock()
		store[key] = value
		mu.Unlock()
		return "+OK\r\n"

	case "GET":
		if len(args) != 2 {
			return "-ERR wrong number of arguments for SET\r\n"
		}
		key := args[1]
		mu.RLock()
		value, ok := store[key]
		mu.RUnlock()
		if !ok {
			return "$-1\r\n"
		}
		return fmt.Sprintf("$%d\r\n%s\r\n", len(value), value)

	default:
		return "-ERR unknown command\r\n"
	}
}
