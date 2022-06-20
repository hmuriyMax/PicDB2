all:
	help

rebuild_user:
	protoc -I api/proto --go_out=pkg --go-grpc_out=require_unimplemented_servers=false:pkg api/proto/user.proto

run_userservice:
	go build ./cmd/userservice/
	sudo ./userservice
	go mod tidy

run_httpserver:
	go build ./cmd/httpserver/
	sudo ./httpserver
	go mod tidy