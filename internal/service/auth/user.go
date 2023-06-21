package auth

import (
	"context"

	authV1 "github.com/a13hander/auth-service-api/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/a13hander/chat-server/internal/domain/model"
	"github.com/a13hander/chat-server/internal/domain/usecase"
)

var _ Client = (*authClient)(nil)
var _ usecase.UserRepo = (*authClient)(nil)

type authClient struct {
	authV1Client authV1.AuthV1Client
}

func NewAuthClient(authV1Client authV1.AuthV1Client) *authClient {
	return &authClient{authV1Client: authV1Client}
}

func (c *authClient) Create(ctx context.Context, u *usecase.CreateUser) (int, error) {
	req := &authV1.CreateRequest{
		User: &authV1.UserInfo{
			Email:    u.Email,
			Username: u.Username,
			Role:     authV1.Role(u.Role),
		},
		Password:        u.Password,
		PasswordConfirm: u.PasswordConfirm,
	}

	cr, err := c.authV1Client.Create(ctx, req, grpc.EmptyCallOption{})
	if err != nil {
		return 0, err
	}

	return int(cr.GetId()), nil
}

func (c *authClient) GetAll(ctx context.Context) ([]*model.User, error) {
	lr, err := c.authV1Client.List(ctx, &emptypb.Empty{}, grpc.EmptyCallOption{})
	if err != nil {
		return nil, err
	}

	users := make([]*model.User, 0, len(lr.GetUser()))
	for _, user := range lr.GetUser() {
		users = append(users, &model.User{
			Id:       int(user.Id),
			Email:    user.GetInfo().GetEmail(),
			Username: user.GetInfo().GetUsername(),
			Role:     int(user.GetInfo().GetRole()),
		})
	}

	return users, nil
}
