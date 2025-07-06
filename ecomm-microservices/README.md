# Ecommerce Microservice

## Tech stacks

- MySQL
- Docker
- Gochi
- gRPC

## GRPC

- https://grpc.io/docs/languages/go/quickstart/#regenerate-grpc-code

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc

export PATH="$PATH:$(go env GOPATH)/bin"

source ~/.zshrc

protoc --proto_path=ecomm-grpc/pb --go_out=ecomm-grpc/pb --go_opt=paths=source_relative \
    --go-grpc_out=ecomm-grpc/pb --go-grpc_opt=paths=source_relative \
    ecomm-grpc/pb/api.proto
```
