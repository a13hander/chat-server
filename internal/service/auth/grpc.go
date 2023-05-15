package auth

import (
	"context"
	"log"

	authV1 "github.com/a13hander/auth-service-api/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/a13hander/chat-server/internal/domain/model"
	"github.com/a13hander/chat-server/internal/domain/usecase"
)

var _ Client = (*grpcClient)(nil)
var _ usecase.UserRepo = (*grpcClient)(nil)

type grpcClient struct {
	grpcClient authV1.AuthV1Client
}

func NewGrpcClient(address string) *grpcClient {
	cc, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalln(err)
	}

	ac := authV1.NewAuthV1Client(cc)

	return &grpcClient{grpcClient: ac}
}

func (c *grpcClient) Create(ctx context.Context, u *usecase.CreateUser) (int, error) {
	req := &authV1.CreateRequest{
		User: &authV1.UserInfo{
			Email:    u.Email,
			Username: u.Username,
			Role:     authV1.Role(u.Role),
		},
		Password:        u.Password,
		PasswordConfirm: u.PasswordConfirm,
	}

	cr, err := c.grpcClient.Create(ctx, req, grpc.EmptyCallOption{})
	if err != nil {
		return 0, err
	}

	return int(cr.GetId()), nil
}

func (c *grpcClient) GetAll(ctx context.Context) ([]*model.User, error) {
	lr, err := c.grpcClient.List(ctx, &emptypb.Empty{}, grpc.EmptyCallOption{})
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
