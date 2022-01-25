package main

import (
	auth_pb "PicDB2/pkg/auth.pb"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net/http"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Not valid method", 405)
	}
	err := r.ParseForm()
	if err != nil {
		return
	}
	x := r.PostForm.Get("login")
	y := r.PostForm.Get("pass")
	conn, err := grpc.Dial(":6000", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	c := auth_pb.NewAuthServerClient(conn)

	res, err := c.GetToken(context.Background(), &auth_pb.LoginData{Login: x, Password: y})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res.GetToken(), res.GetIsAuthorised())
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Path[1:]
	str = "./cmd/httpserver/html/login.html"
	println(str)
	var tpl = template.Must(template.ParseFiles(str))
	err := tpl.Execute(w, nil)
	if err != nil {
		return
	}
}
