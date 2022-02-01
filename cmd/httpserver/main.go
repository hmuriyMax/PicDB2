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
	err := tpl.Execute(w, nil)
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
		mess = Errpass
	}
	params := map[string]string{
		"login":  r.URL.Query().Get("login"),
		"status": string(template.HTML(mess)),
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
	err := tpl.Execute(w, nil)
	if err != nil {
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fileServer := http.FileServer(http.Dir("./cmd/httpserver/html/res/"))
	http.Handle("/res/", http.StripPrefix("/res", fileServer))
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/reg", regHandler)
	http.HandleFunc("/newser", newUserHandler)
	http.HandleFunc("/logout", loutHandler)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return
	}
}
