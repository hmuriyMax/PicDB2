package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func CheckIn(arr []string, s string) bool {
	for i := 0; i < len(arr); i++ {
		if arr[i] == s[:len(arr[i])] {
			return true
		}
	}
	return false
}

var endf = []string{".html", ".css", ".js"}
var maxl = 5

func indexHandler(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Path[1:]
	if len(str) < 5 || CheckIn(endf, str[:maxl]) {
		str += "./cmd/httpserver/html/index.html"
	}
	ip := r.RemoteAddr
	log.Printf("IP %s GET %s", ip, str)
	var tpl = template.Must(template.ParseFiles(str))
	var params map[string]template.HTML
	_, err := r.Cookie("user_id")
	if err == nil {
		params = map[string]template.HTML{
			"buttons": template.HTML(loggedButtons),
		}
	} else {
		params = map[string]template.HTML{
			"buttons": template.HTML(unLoggedButtons),
		}
	}
	pic, err := r.Cookie("profile_picture_url")
	if err == nil {
		params["profpic"] = template.HTML(pic.Value)
	} else {
		params["profpic"] = "/res/img/default_profile_pic.png"
	}
	err = tpl.Execute(w, params)
	if err != nil {
		return
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Path[1:]
	str = "./cmd/httpserver/html/login.html"
	ip := r.RemoteAddr
	log.Printf("IP %s GET %s", ip, str)
	var tpl = template.Must(template.ParseFiles(str))
	mess := ""
	if r.URL.Query().Get("status") == "unath" {
		mess = errPass
	}
	params := map[string]string{
		"login":  r.URL.Query().Get("login"),
		"status": mess,
	}
	err := tpl.Execute(w, params)
	if err != nil {
		return
	}
}

func regHandler(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Path[1:]
	str = "./cmd/httpserver/html/register.html"
	ip := r.RemoteAddr
	log.Printf("IP %s GET %s", ip, str)
	var tpl = template.Must(template.ParseFiles(str))
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	uname := r.URL.Query().Get("uname")
	status := r.URL.Query().Get("status")
	if status == "unalreadyexist" {
		status = errUsername
	} else if status == "emalreadyexist" {
		status = errEmail
	}
	params := map[string]string{
		"name":   name,
		"email":  email,
		"uname":  uname,
		"status": status,
	}
	err := tpl.Execute(w, params)
	if err != nil {
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./cmd/httpserver/html/res/"))
	mux.Handle("/res/", http.StripPrefix("/res", fileServer))
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/auth", authHandler)
	mux.HandleFunc("/login", loginHandler)
	mux.HandleFunc("/reg", regHandler)
	mux.HandleFunc("/newser", newUserHandler)
	mux.HandleFunc("/logout", loutHandler)
	mux.HandleFunc("/profile", profileHandler)
	log.Printf("HTTP-server started! Port: %s", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
