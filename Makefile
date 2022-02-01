all:
	help

#rebuild_auth:
#	protoc -I api/proto --go_out=pkg/auth.pb --go-grpc_out=require_unimplemented_servers=false:pkg api/proto/auth.proto

rebuild_user:
	protoc -I api/proto --go_out=pkg/user.pb --go-grpc_out=require_unimplemented_servers=false:pkg api/proto/user.proto

#run_auth_server:
#	go mod tidy
#	go build ./cmd/authservice/
#	./authservice

run_user_server:
	go mod tidy
	go build ./cmd/userservice/
	./userservice

run_http_server:
	go mod tidy
	go build ./cmd/httpserver/
	./httpserver