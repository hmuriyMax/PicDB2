all:
	help

rebuild:
	protoc -I api/proto --go_out=pkg/auth.pb --go-grpc_out=require_unimplemented_servers=false:pkg api/proto/auth.proto

run_auth_server:
	go build ./cmd/authservice/
	./authservice

run_http_server:
	go build ./cmd/httpserver/
	./httpserver