package main

// Go + gRPC
func main() {
	grpcServer := NewGRPCServer(":9000")
	grpcServer.Run()
}
