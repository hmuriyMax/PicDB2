package main

import (
	userPB "PicDB2/pkg/user.pb"
	"log"
	"net/http"
	"strconv"
	"time"
)

var errPass = "Неверный логин или пароль!"
var errUsername = "Пользователь с таким именем пользователя уже существует!"
var errEmail = "Пользователь с такой почтой уже существует!"
var unLoggedButtons = " <a href=\"/login\" id=\"LogInBut\">Войти</a>\n <a href=\"/reg\" id=\"RegBut\">Регистрация</a>"
var loggedButtons = " <a href=\"/profile\" id=\"ProfBut\">Профиль</a>\n <a href=\"/logout\" id=\"LogOutBut\">Выйти</a>"

func SetCookie(w http.ResponseWriter, cookie http.Cookie) {
	http.SetCookie(w, &cookie)
	log.Printf("COOKIE is set: %s: %s expires %s", cookie.Name, cookie.Value, cookie.Expires)
}

func Redirect(w http.ResponseWriter, url string, status int) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, req, url, status)
}

func SetTokenCookies(res *userPB.LoginStatus, w http.ResponseWriter) bool {
	if res.IsAuthorised {
		parse, err := time.Parse(time.ANSIC, res.GetToken().GetExpires())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
		passCookie := http.Cookie{
			Name:       "passhash",
			Value:      res.GetToken().GetToken(),
			Path:       "/",
			Domain:     "",
			Expires:    parse,
			RawExpires: parse.Format(time.UnixDate),
			MaxAge:     60 * 60 * 24 * 60,
			Secure:     true,
			HttpOnly:   true,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		}
		idCookie := http.Cookie{
			Name:       "user_id",
			Value:      strconv.Itoa(int(res.GetToken().GetUid())),
			Path:       "/",
			Domain:     "",
			Expires:    parse,
			RawExpires: parse.Format(time.UnixDate),
			MaxAge:     60 * 60 * 24 * 60,
			Secure:     true,
			HttpOnly:   true,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		}
		SetCookie(w, passCookie)
		SetCookie(w, idCookie)
		return true
	} else {
		http.Error(w, "Server error: not authorised", http.StatusInternalServerError)
		log.Fatal("Server error: not authorised")
		return false
	}
}

func DelCookie(w http.ResponseWriter, cookieName string) {
	c := http.Cookie{
		Name:   cookieName,
		MaxAge: -1,
	}
	http.SetCookie(w, &c)
}
