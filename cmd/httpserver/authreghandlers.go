package main

import (
	auth_pb "PicDB2/pkg/auth.pb"
	"context"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

func SetCookie(w http.ResponseWriter, cookie http.Cookie) {
	http.SetCookie(w, &cookie)
	log.Printf("COOKIE is set: %s: %s expires %s", cookie.Name, cookie.Value, cookie.Expires)
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(":6000", grpc.WithInsecure())
	c := auth_pb.NewAuthServerClient(conn)

	ph, err := r.Cookie("passhash")
	uh, err := r.Cookie("user_id")
	if ph.Value != "" && uh.Value != "" {
		cookieUserId, err := strconv.Atoi(uh.Value)
		authorised, err := c.IsAuthorised(context.Background(), &auth_pb.Token{Token: ph.Value, Uid: int32(cookieUserId), Expires: ph.RawExpires})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Fatal(err)
		}
		if authorised.IsAuthorised {
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				log.Fatal(err)
			}
			http.Redirect(w, req, "/", http.StatusAccepted)
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
	res, err := c.GetToken(context.Background(), &auth_pb.LoginData{Login: x, Password: y})
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
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	str := r.URL.Path[1:]
	str = "./cmd/httpserver/html/login.html"
	ip := r.RemoteAddr
	log.Printf("IP %s GET %s", ip, str)
	var tpl = template.Must(template.ParseFiles(str))
	err := tpl.Execute(w, nil)
	if err != nil {
		return
	}
}
