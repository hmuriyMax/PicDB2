package main

import (
	auth_pb "PicDB2/pkg/auth.pb"
	"context"
	"flag"
	"google.golang.org/grpc"
	"log"
)

func main() {
	flag.Parse()
	if flag.NArg() < 2 {
		log.Fatal("not enough args")
	}

	x := flag.Arg(0)
	y := flag.Arg(1)

	conn, err := grpc.Dial(":6000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := auth_pb.NewAuthServerClient(conn)

	res, err := c.GetToken(context.Background(), &auth_pb.LoginData{Login: x, Password: y})
	if err != nil {
		log.Fatal(err)
	}

	log.Println(res.GetToken(), res.GetIsAuthorised())
}
