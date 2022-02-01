package main

import (
	userPB "PicDB2/pkg/user.pb"
	"context"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
	"time"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(":6000", grpc.WithInsecure())
	c := userPB.NewUserServerClient(conn)

	ph, err := r.Cookie("passhash")
	uh, err := r.Cookie("user_id")
	if ph.Value != "" && uh.Value != "" {
		cookieUserId, err := strconv.Atoi(uh.Value)
		authorised, err := c.IsAuthorised(context.Background(), &userPB.Token{Token: ph.Value, Uid: int32(cookieUserId), Expires: ph.RawExpires})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
		if authorised.IsAuthorised {
			Redirector(w, "/", http.StatusFound)
			return
		}
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Not valid method", http.StatusMethodNotAllowed)
		log.Fatal(err)
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	x := r.PostForm.Get("login")
	y := r.PostForm.Get("pass")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	res, err := c.GetToken(context.Background(), &userPB.LoginData{Login: x, Password: y})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
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
		Redirector(w, "/", http.StatusFound)
	}
	Redirector(w, "/login?login="+x+"&status=unath", http.StatusSeeOther)
}

func newUserHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(":6000", grpc.WithInsecure())
	c := userPB.NewUserServerClient(conn)
	//новый GRPC сервер для сохранения информации

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
	user, err := c.NewUser(context.Background(), &userPB.LoginData{Login: "", Password: ""})
	if err != nil {
		return
	}
	user.GetToken()
}

func loutHandler(writer http.ResponseWriter, request *http.Request) {

}
