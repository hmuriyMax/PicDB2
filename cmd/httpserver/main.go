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
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		return
	}
}
