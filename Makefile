THIS_FILE := $(lastword $(MAKEFILE_LIST))

proto:
	 mkdir -p out && protoc \
	./pb/model.proto ./pb/customer.proto ./pb/transaction.proto ./pb/transactionItem.proto ./pb/profile.proto ./pb/lang.proto \
  --go_out=:. --go-grpc_out=:. \
	--go_opt=Mpb/transaction.proto=github.com/bdarge/api/out/transaction \
	--go_opt=Mpb/transactionItem.proto=github.com/bdarge/api/out/transactionItem \
	--go_opt=Mpb/customer.proto=github.com/bdarge/api/out/customer \
	--go_opt=Mpb/model.proto=github.com/bdarge/api/out/model \
	--go_opt=Mpb/profile.proto=github.com/bdarge/api/out/profile \
	--go_opt=Mpb/lang.proto=github.com/bdarge/api/out/lang \
	--go_opt=module=github.com/bdarge/api \
	--go-grpc_opt=Mpb/transaction.proto=github.com/bdarge/api/out/transaction \
	--go-grpc_opt=Mpb/transactionItem.proto=github.com/bdarge/api/out/transactionItem \
	--go-grpc_opt=Mpb/customer.proto=github.com/bdarge/api/out/customer \
	--go-grpc_opt=Mpb/model.proto=github.com/bdarge/api/out/model \
	--go-grpc_opt=Mpb/profile.proto=github.com/bdarge/api/out/profile \
	--go-grpc_opt=Mpb/lang.proto=github.com/bdarge/api/out/lang \
	--go-grpc_opt=module=github.com/bdarge/api


server:
	go run cmd/main.go

build:
	@$(MAKE) -f $(THIS_FILE) proto; docker build -f Dockerfile.grpc -t api --target dev . --no-cache
