package main

import (
	userPB "PicDB2/pkg/user_pb"
	"bytes"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var errPass = "Неверный логин или пароль!"
var errUsername = "Пользователь с таким именем пользователя уже существует!"
var errEmail = "Пользователь с такой почтой уже существует!"
var unLoggedButtons = " <a href=\"/login\" id=\"LogInBut\">Войти</a>\n <a href=\"/reg\" id=\"RegBut\">Регистрация</a>"
var loggedButtons = " <a href=\"/profile\" id=\"ProfBut\">Профиль</a>\n <a href=\"/logout\" id=\"LogOutBut\">Выйти</a>"
var HTMLpath = "./cmd/httpserver/html/"
var authAge = 60 * 60 * 24 * 60

func SetCookie(w http.ResponseWriter, name, value string, MaxAge int) {
	tmp := http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: MaxAge,
		Domain: "/",
	}
	HTTPSetCookie(w, tmp)
}

func HTTPSetCookie(w http.ResponseWriter, cookie http.Cookie) {
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
		//parse, err := time.Parse(time.ANSIC, res.GetToken().GetExpires())
		//if err != nil {
		//	http.Error(w, err.Error(), http.StatusInternalServerError)
		//	log.Fatal(err)
		//}
		SetCookie(w, "passhash", res.GetToken().GetToken(), 60*60*24*60)
		SetCookie(w, "user_id", strconv.Itoa(int(res.GetToken().GetUid())), 60*60*24*60)
		//passCookie := http.Cookie{
		//	Name:       "passhash",
		//	Value:      res.GetToken().GetToken(),
		//	Path:       "/",
		//	Domain:     "",
		//	Expires:    parse,
		//	RawExpires: parse.Format(time.UnixDate),
		//	MaxAge:     60 * 60 * 24 * 60,
		//	Secure:     true,
		//	HttpOnly:   true,
		//	SameSite:   0,
		//	Raw:        "",
		//	Unparsed:   nil,
		//}
		//idCookie := http.Cookie{
		//	Name:       "user_id",
		//	Value:      strconv.Itoa(int(res.GetToken().GetUid())),
		//	Path:       "/",
		//	Domain:     "",
		//	Expires:    parse,
		//	RawExpires: parse.Format(time.UnixDate),
		//	MaxAge:     60 * 60 * 24 * 60,
		//	Secure:     true,
		//	HttpOnly:   true,
		//	SameSite:   0,
		//	Raw:        "",
		//	Unparsed:   nil,
		//}
		//HTTPSetCookie(w, passCookie)
		//HTTPSetCookie(w, idCookie)
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

func GenerateHeader(r *http.Request) (string, error) {
	headerPath := HTMLpath + "header.html"
	var headerTmp = template.Must(template.ParseFiles(headerPath))
	bt := unLoggedButtons
	_, err := r.Cookie("user_id")
	if err == nil {
		bt = loggedButtons
	}
	params := map[string]template.HTML{
		"buttons": template.HTML(bt),
	}
	var res bytes.Buffer
	err = headerTmp.Execute(&res, params)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

func AddCookies(params *map[string]template.HTML, r *http.Request) {
	for _, cookie := range r.Cookies() {
		(*params)[cookie.Name] = template.HTML(cookie.Value)
	}
}

func DialUserService() *userPB.UserServerClient {
	conn, err := grpc.Dial(":6000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {

	}
	c := userPB.NewUserServerClient(conn)
	return &c
}

func GetCookie(r *http.Request, name string) string {
	cookie, err := r.Cookie(name)
	if err != nil {
		log.Println(err)
		return ""
	}
	return cookie.Value
}
