package main

import (
	"context"
	userPB "github.com/hmuriyMax/PicDB2/pkg/user_pb"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func profileHandler(w http.ResponseWriter, r *http.Request) {
	profilePath := HTMLpath + "profile.html"
	ip := r.RemoteAddr
	log.Printf("IP %s GET %s", ip, profilePath)
	var profileTpl = template.Must(template.ParseFiles(profilePath))

	header, err := GenerateHeader(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := map[string]template.HTML{
		"header": template.HTML(header),
	}
	//TODO: добавить вкладки с подписчиками, подписками и редактированием профиля через GET
	//TODO: Возможность заходить на чужой профиль (через /profile?id=01545 или /profile01545)
	//TODO: ошибка при неавт. переходе в страницу профиля
	c := *DialUserService()
	uID, err := strconv.Atoi(GetCookie(r, "user_id"))
	if err != nil {
		log.Print(err)
	}
	data, err := c.GetPartUserData(context.Background(), &userPB.UserId{Id: int32(uID)})
	if err != nil {
		return
	}
	AddCookies(&params, r)
	params["name"] = template.HTML(data.Name)
	params["follnum"] = template.HTML("0")
	params["subsnum"] = template.HTML("0")
	params["profile_picture_url"] = template.HTML(data.GetProfilePicURL())
	err = profileTpl.Execute(w, params)
	if err != nil {
		return
	}
}
