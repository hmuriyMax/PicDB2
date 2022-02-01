package main

import (
	"PicDB2/pkg/user"
	api "PicDB2/pkg/user.pb"
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