package auth

import (
	api "PicDB2/pkg/auth.pb"
	"context"
)

type GRPCServer struct{}

func (s *GRPCServer) GetToken(ctx context.Context, logdt *api.LoginData) (*api.LoginStatus, error) {
	userlog := logdt.GetLogin()
	pass := logdt.GetPassword()
	status, err := CheckUser(userlog, pass)

	var lstat = api.LoginStatus{}
	if err != nil {
		return nil, err
	}
	if status < 0 {
		lstat.Token = nil
		lstat.IsAuthorised = false
	} else {
		token, err := GetToken(status)
		if err != nil {
			return nil, err
		}
		lstat.Token = &api.Token{Token: token}
		lstat.IsAuthorised = true
	}

	return &lstat, nil
}

func (s *GRPCServer) IsAuthorised(ctx context.Context, tok *api.Token) (*api.LoginStatus, error) {
	status, err := CheckToken(tok.Token)
	if err != nil {
		return nil, err
	}
	return &api.LoginStatus{
		Token:        tok,
		IsAuthorised: status}, nil
}

func (s *GRPCServer) NewUser(ctx context.Context, logdt *api.LoginData) (*api.LoginStatus, error) {
	userlog := logdt.GetLogin()
	pass := logdt.GetPassword()
	id, err := InsertUser(userlog, pass)
	if err != nil {
		return nil, err
	}
	token, err := GetToken(id)
	if err != nil {
		return nil, err
	}
	return &api.LoginStatus{
		Token:        &api.Token{Token: token},
		IsAuthorised: true}, nil
}

func (s *GRPCServer) mustEmbedUnimplementedAuthServerServer() {
	return
}