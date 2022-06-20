package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	indexPath := HTMLpath + "index.html"
	ip := r.RemoteAddr
	log.Printf("IP %s GET %s", ip, indexPath)
	var indexTmp = template.Must(template.ParseFiles(indexPath))
	header, err := GenerateHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := map[string]template.HTML{
		"header": template.HTML(header),
	}
	pic, err := r.Cookie("profile_picture_url")
	if err == nil {
		params["profile_picture_url"] = template.HTML(pic.Value)
	} else {
		params["profile_picture_url"] = "/res/img/default_profile_pic.png"
	}
	err = indexTmp.Execute(w, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
	regPath := HTMLpath + "register.html"
	ip := r.RemoteAddr
	log.Printf("IP %s GET %s", ip, regPath)
	var regTpl = template.Must(template.ParseFiles(regPath))

	header, err := GenerateHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := map[string]template.HTML{
		"header": template.HTML(header),
	}

	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	uname := r.URL.Query().Get("uname")
	status := r.URL.Query().Get("status")

	if status == "unalreadyexist" {
		status = errUsername
	} else if status == "emalreadyexist" {
		status = errEmail
	}

	params["name"] = template.HTML(name)
	params["email"] = template.HTML(email)
	params["uname"] = template.HTML(uname)
	params["status"] = template.HTML(status)
	err = regTpl.Execute(w, params)
	if err != nil {
		return
	}
}

func main() {
	port := "80"
	if len(os.Args) > 1 {
		port = os.Args[1]
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
	log.Printf("HTTP-server started! http://localhost:%s", port)
	err := http.ListenAndServe(":"+port, mux)
	if err != nil {
		log.Fatal(err)
	}
}
