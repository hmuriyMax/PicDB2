all:
	help

#rebuild_auth:
#	protoc -I api/proto --go_out=pkg/auth.pb --go-grpc_out=require_unimplemented_servers=false:pkg api/proto/auth.proto

rebuild_user:
	protoc -I api/proto --go_out=pkg --go-grpc_out=require_unimplemented_servers=false:pkg api/proto/user.proto

#run_auth_server:
#	go build ./cmd/authservice/
#	./authservice
#	go mod tidy

run_user_server:
	go build ./cmd/userservice/
	./userservice
	go mod tidy

run_http_server:
	go build ./cmd/httpserver/
	./httpserver
	go mod tidy