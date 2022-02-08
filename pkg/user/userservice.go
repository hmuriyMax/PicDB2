package user

import (
	api "PicDB2/pkg/user.pb"
	"context"
	"fmt"
)

type UserData struct {
	userId        int32
	name          string
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
	udata := UserData{
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
	err := DBDeleteUser(id.GetId())
	if err != nil {
		return nil, err
	}
	return &api.Status{Code: 0}, nil
}

func (s *GRPCServer) GetUserData(ctx context.Context, id *api.UserId) (*api.UserData, error) {
	data, err := DBGetUserData(id.GetId())
	if err != nil {
		return nil, err
	}
	return &api.UserData{
		Id:            data.userId,
		Name:          data.name,
		Gender:        data.gender,
		Bday:          data.birthday,
		ProfilePicURL: data.profilePicURL,
		Unqhash:       data.uniqueKey,
		Email:         data.email,
	}, nil
}
