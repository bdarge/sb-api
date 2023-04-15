proto:
	 mkdir -p out && protoc ./pb/model.proto ./pb/customer.proto ./pb/transaction.proto  --go_out=:. --go-grpc_out=:. \
	--go_opt=Mpb/transaction.proto=github.com/bdarge/api/out/transaction \
	--go_opt=Mpb/customer.proto=github.com/bdarge/api/out/customer \
	--go_opt=Mpb/model.proto=github.com/bdarge/api/out/model \
	--go_opt=module=github.com/bdarge/api \
	--go-grpc_opt=Mpb/transaction.proto=github.com/bdarge/api/out/transaction \
	--go-grpc_opt=Mpb/customer.proto=github.com/bdarge/api/out/customer \
	--go-grpc_opt=Mpb/model.proto=github.com/bdarge/api/out/model \
	--go-grpc_opt=module=github.com/bdarge/api

server:
	go run cmd/main.go