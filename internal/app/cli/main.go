package cli

import (
	"context"
	"fmt"

	"github.com/a13hander/chat-server/internal/domain/usecase"
)

type AppCli interface {
	Run(ctx context.Context) error
}

type appCli struct {
	createUserUseCase usecase.CreateUserUseCase
	userListUseCase   usecase.ListUserUseCase
}

func NewAppCli(createUserUseCase usecase.CreateUserUseCase, userListUseCase usecase.ListUserUseCase) *appCli {
	return &appCli{createUserUseCase: createUserUseCase, userListUseCase: userListUseCase}
}

func (c *appCli) Run(ctx context.Context) error {
	req := &usecase.CreateUserRequest{}

	fmt.Println("Enter user info:")
	fmt.Print("email:")
	_, _ = fmt.Scanf("%s\n", &req.Email)
	fmt.Print("username:")
	_, _ = fmt.Scanf("%s\n", &req.Username)
	fmt.Print("password:")
	_, _ = fmt.Scanf("%s\n", &req.Password)
	fmt.Print("password confirm:")
	_, _ = fmt.Scanf("%s\n", &req.PasswordConfirm)
	fmt.Print("role:")
	_, _ = fmt.Scanf("%s\n", &req.Role)

	id, err := c.createUserUseCase.Run(ctx, req)
	if err != nil {
		return err
	}

	fmt.Printf("New user id %d\n\n", id)

	users, err := c.userListUseCase.Run(ctx)
	if err != nil {
		return err
	}

	for _, u := range users {
		fmt.Printf("%v\n", u)
	}

	return nil
}
