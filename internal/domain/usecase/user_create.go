package usecase

import (
	"context"
	"fmt"

	"github.com/a13hander/chat-server/internal/domain/util"
)

type CreateUserRequest struct {
	Email           string
	Username        string
	Password        string
	PasswordConfirm string
	Role            int
}

func (r *CreateUserRequest) String() string {
	return fmt.Sprintf("%v", *r)
}

type createUserUseCase struct {
	validator UserValidator
	repo      UserRepo
	l         util.Logger
}

func NewCreateUserUseCase(validator UserValidator, repo UserRepo, l util.Logger) *createUserUseCase {
	return &createUserUseCase{validator: validator, repo: repo, l: l}
}

func (c *createUserUseCase) Run(ctx context.Context, req *CreateUserRequest) (int, error) {
	err := c.validator.ValidateCreating(req)
	if err != nil {
		return 0, err
	}

	u := CreateUser{
		Email:           req.Email,
		Username:        req.Username,
		Password:        req.Password,
		PasswordConfirm: req.PasswordConfirm,
		Role:            req.Role,
	}

	id, err := c.repo.Create(ctx, &u)
	if err != nil {
		c.l.ErrorWithCtx("не удалось создать пользователя", map[string]any{
			"err":     err.Error(),
			"payload": req.String(),
		})
		return 0, err
	}

	return id, nil
}
