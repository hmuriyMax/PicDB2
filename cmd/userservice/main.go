package main

import (
	"github.com/hmuriyMax/PicDB2/pkg/user"
	api "github.com/hmuriyMax/PicDB2/pkg/user_pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	srv := &user.GRPCServer{}
	api.RegisterUserServerServer(s, srv)
	l, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}
	err = s.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
}
