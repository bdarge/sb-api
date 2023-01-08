proto:
	protoc pkg/**/customer.proto --go_out=:. --go-grpc_out=:.
	protoc pkg/**/disposition.proto --go_out=:. --go-grpc_out=:.

server:
	go run cmd/main.go