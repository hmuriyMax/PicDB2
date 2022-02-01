package user

import (
	api "PicDB2/pkg/user.pb"
	"context"
	"fmt"
	"time"
)

type UserData struct {
	userId        int32
	username      string
	birthday      time.Time
	gender        string
	profilePicURL string
	uniqueKey     string
}

func (s *GRPCServer) UpdateUser(ctx context.Context, data *api.UserData) (*api.Status, error) {
	defer fmt.Printf("\n")
	ttime, err := time.Parse(time.RFC3339, data.GetBday())
	if err != nil {
		return nil, err
	}
	udata := UserData{userId: data.GetId(),
		username:      data.GetName(),
		birthday:      ttime,
		gender:        data.GetGender(),
		profilePicURL: data.GetProfilePicURL(),
		uniqueKey:     data.GetUnqhash(),
	}
	print(udata.username)
	return nil, nil
}

func (s *GRPCServer) DeleteUser(ctx context.Context, id *api.UserId) (*api.Status, error) {
	return nil, nil
}

func (s *GRPCServer) GetUserData(ctx context.Context, id *api.UserId) (*api.UserData, error) {
	return nil, nil
}
