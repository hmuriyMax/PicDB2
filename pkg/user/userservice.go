package user

import (
	"context"
	"fmt"
	api "github.com/hmuriyMax/PicDB2/pkg/user_pb"
	"strings"
)

type Data struct {
	userId        int32
	name          string
	username      string
	email         string
	birthday      string
	gender        string
	profilePicURL string
	uniqueKey     string
}

func (s *GRPCServer) UpdateUser(ctx context.Context, data *api.UserData) (*api.Status, error) {
	defer fmt.Printf("\n")
	//ttime, err := time.Parse(time.RFC3339, data.GetBday())
	//if err != nil {
	//	return nil, err
	//}
	udata := Data{
		userId:        data.GetId(),
		name:          data.GetName(),
		birthday:      data.GetBday(),
		email:         data.GetEmail(),
		gender:        data.GetGender(),
		profilePicURL: data.GetProfilePicURL(),
		uniqueKey:     data.GetUnqhash(),
	}
	err := DBUpdateUserData(&udata)
	if err != nil {
		return nil, err
	}
	return &api.Status{Code: 0}, nil
}

func (s *GRPCServer) DeleteUser(ctx context.Context, id *api.UserId) (*api.Status, error) {
	defer fmt.Printf("\n")
	err := DBDeleteUser(id.GetId())
	if err != nil {
		return nil, err
	}
	return &api.Status{Code: 0}, nil
}

func (s *GRPCServer) GetFullUserData(ctx context.Context, id *api.UserId) (*api.UserData, error) {
	defer fmt.Printf("\n")
	data, err := DBGetFullUserData(id.GetId())
	if err != nil {
		return nil, err
	}
	return &api.UserData{
		Id:            data.userId,
		Name:          strings.TrimSpace(data.name),
		Uname:         strings.TrimSpace(data.username),
		Gender:        strings.TrimSpace(data.gender),
		Bday:          strings.TrimSpace(data.birthday),
		ProfilePicURL: strings.TrimSpace(data.profilePicURL),
		Unqhash:       strings.TrimSpace(data.uniqueKey),
		Email:         strings.TrimSpace(data.email),
	}, nil
}

func (s *GRPCServer) GetPartUserData(ctx context.Context, id *api.UserId) (*api.UserDataS, error) {
	defer fmt.Printf("\n")
	data, err := DBGetShortUserData(id.GetId())
	if err != nil {
		return nil, err
	}
	return &api.UserDataS{
		Id:            data.userId,
		Name:          strings.TrimSpace(data.name),
		Username:      strings.TrimSpace(data.username),
		ProfilePicURL: strings.TrimSpace(data.profilePicURL),
	}, nil
}
