package usecase

import (
	"context"

	"github.com/a13hander/chat-server/internal/domain/model"
)

type CreateUser struct {
	Email           string
	Username        string
	Password        string
	PasswordConfirm string
	Role            int
}

type UserValidator interface {
	ValidateCreating(r *CreateUserRequest) error
}

type UserRepo interface {
	Create(ctx context.Context, user *CreateUser) (int, error)
	GetAll(ctx context.Context) ([]*model.User, error)
}

type CreateUserUseCase interface {
	Run(ctx context.Context, req *CreateUserRequest) (int, error)
}

type ListUserUseCase interface {
	Run(ctx context.Context) ([]*model.User, error)
}
