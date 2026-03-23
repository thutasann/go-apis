package main

// Go + gRPC
func main() {
	httpServer := NewHttpServer(":8000")
	httpServer.Run()

	grpcServer := NewGRPCServer(":9000")
	grpcServer.Run()
}
