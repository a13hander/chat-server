package auth

import (
	"context"

	"github.com/a13hander/chat-server/internal/domain/model"
	"github.com/a13hander/chat-server/internal/domain/usecase"
)

type Client interface {
	Create(ctx context.Context, u *usecase.CreateUser) (int, error)
	GetAll(ctx context.Context) ([]*model.User, error)
}
