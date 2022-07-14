package clientutil

import (
	"context"
	"github.com/go-funcards/funapi/internal/gin/httputil"
	v1User "github.com/go-funcards/funapi/proto/user_service/v1"
)

func UsersRequest(id string) *v1User.UsersRequest {
	return &v1User.UsersRequest{
		PageIndex: 0,
		PageSize:  1,
		UserIds:   []string{id},
	}
}

func GetUser(ctx context.Context, client v1User.UserClient, id string) (*v1User.UserResponse, error) {
	response, err := client.GetUsers(ctx, UsersRequest(id))
	if err != nil {
		return nil, err
	}
	if len(response.GetUsers()) != 1 {
		return nil, httputil.ErrNotFound
	}
	return response.GetUsers()[0], nil
}
