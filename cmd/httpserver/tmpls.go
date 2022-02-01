package main

import (
	"log"
	"net/http"
)

var Errpass = "Неверный логин или пароль!"

func SetCookie(w http.ResponseWriter, cookie http.Cookie) {
	http.SetCookie(w, &cookie)
	log.Printf("COOKIE is set: %s: %s expires %s", cookie.Name, cookie.Value, cookie.Expires)
}

func Redirector(w http.ResponseWriter, url string, status int) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, req, url, status)
}
