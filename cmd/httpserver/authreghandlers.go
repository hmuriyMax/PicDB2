package main

import (
	userPB "PicDB2/pkg/user.pb"
	"context"
	"google.golang.org/grpc"
	"log"
	"net/http"
)

func authHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(":6000")
	c := userPB.NewUserServerClient(conn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}

	_, err = r.Cookie("passhash")
	_, err2 := r.Cookie("user_id")
	if err == nil && err2 == nil {
		DelCookie(w, "user_id")
		DelCookie(w, "passhash")
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Not valid method", http.StatusMethodNotAllowed)
		log.Println(err)
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
	x := r.PostForm.Get("login")
	y := r.PostForm.Get("pass")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
	res, err := c.GetToken(context.Background(), &userPB.LoginData{Login: x, Password: y})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
	if res.IsAuthorised {
		SetTokenCookies(res, w)
		// TODO: Сохранять в куки имя пользователя и URL картинки
		Redirect(w, "/", http.StatusFound)
	} else {
		Redirect(w, "/login?login="+x+"&status=unath", http.StatusSeeOther)
	}
}

func newUserHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := grpc.Dial(":6000")
	c := userPB.NewUserServerClient(conn)
	if r.Method != http.MethodPost {
		http.Error(w, "Not valid method", http.StatusMethodNotAllowed)
		log.Println(err)
	}
	err = r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
	postedName := r.PostForm.Get("name")
	postedEmail := r.PostForm.Get("email")
	postedUname := r.PostForm.Get("uname")
	postedPass := r.PostForm.Get("pass")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
	userid, err := c.NewUser(context.Background(), &userPB.LoginData{Login: postedUname, Password: postedPass})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
	if userid.GetId() == -1 {
		Redirect(w, "/reg?name="+postedName+"&email="+
			postedEmail+"&uname="+postedUname+"&status=unalreadyexist", http.StatusSeeOther)
		return
	}
	lstat, err := c.GetToken(context.Background(), &userPB.LoginData{Login: postedUname, Password: postedPass})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println(err)
	}
	_, err = c.UpdateUser(context.Background(), &userPB.UserData{Id: userid.GetId(), Name: postedName, Email: postedEmail})
	if err != nil {
		_, err := c.DeleteUser(context.Background(), &userPB.UserId{Id: userid.GetId()})
		if err != nil {
			return
		}
		Redirect(w, "/reg?name="+postedName+"&email="+
			postedEmail+"&uname="+postedUname+"&status=emalreadyexist", http.StatusSeeOther)
		return
	}
	SetTokenCookies(lstat, w)
	Redirect(w, "/", http.StatusFound)
}

func loutHandler(w http.ResponseWriter, _ *http.Request) {
	DelCookie(w, "user_id")
	DelCookie(w, "passhash")
	Redirect(w, "/", http.StatusFound)
}
