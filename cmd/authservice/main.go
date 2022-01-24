package main

import (
	"PicDB2/pkg/auth"
	api "PicDB2/pkg/auth.pb"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	s := grpc.NewServer()
	srv := &auth.GRPCServer{}
	api.RegisterAuthServerServer(s, srv)
	l, err := net.Listen("tcp", ":6000")
	if err != nil {
		log.Fatal(err)
	}
	err = s.Serve(l)
	if err != nil {
		log.Fatal(err)
	}
}
