package user

import (
	api "PicDB2/pkg/user.pb"
	"context"
	"fmt"
	"log"
)

type GRPCServer struct{}

func (s *GRPCServer) GetToken(ctx context.Context, logdt *api.LoginData) (*api.LoginStatus, error) {
	defer fmt.Printf("\n")
	userlog := logdt.GetLogin()
	pass := logdt.GetPassword()
	status, err := DBCheckUser(userlog, pass)

	var lstat = api.LoginStatus{}
	if err != nil {
		return nil, err
	}
	if status < 0 {
		lstat.Token = nil
		lstat.IsAuthorised = false
	} else {
		token, err := DBGetToken(status)
		if err != nil {
			return nil, err
		}
		lstat.Token = &api.Token{Token: token.Token, Uid: int32(token.Userid), Expires: token.Expires}
		lstat.IsAuthorised = true
	}
	return &lstat, nil
}

func (s *GRPCServer) IsAuthorised(ctx context.Context, tok *api.Token) (*api.LoginStatus, error) {
	defer fmt.Printf("\n")
	status, err := DBCheckToken(tok.Token)
	if err != nil {
		return nil, err
	}
	return &api.LoginStatus{
		Token:        tok,
		IsAuthorised: status}, nil
}

func (s *GRPCServer) NewUser(ctx context.Context, logdt *api.LoginData) (*api.UserId, error) {
	defer fmt.Printf("\n")
	userlog := logdt.GetLogin()
	pass := logdt.GetPassword()
	id, err := DBInsertUser(userlog, pass)
	if err != nil {
		log.Print(err)
		return nil, err
	}
	return &api.UserId{Id: int32(id)}, nil
}

func (s *GRPCServer) mustEmbedUnimplementedAuthServerServer() {
	return
}
