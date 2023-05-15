package app

import (
	"context"

	"github.com/a13hander/chat-server/internal/config"
	"github.com/a13hander/chat-server/internal/domain/usecase"
	"github.com/a13hander/chat-server/internal/domain/util"
	"github.com/a13hander/chat-server/internal/domain/validator"
	"github.com/a13hander/chat-server/internal/service/auth"
	"github.com/a13hander/chat-server/internal/service/logger"
)

type serviceProvider struct {
	logger util.Logger

	repo struct {
		userRepo usecase.UserRepo
	}

	validator struct {
		userValidator usecase.UserValidator
	}

	useCase struct {
		createUserUseCase usecase.CreateUserUseCase
		userListUseCase   usecase.ListUserUseCase
	}
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (c *serviceProvider) GetLogger(_ context.Context) util.Logger {
	if c.logger == nil {
		c.logger = logger.NewLogger()
	}

	return c.logger
}

func (c *serviceProvider) GetUserRepo(_ context.Context) usecase.UserRepo {
	if c.repo.userRepo == nil {
		c.repo.userRepo = auth.NewGrpcClient(config.GetConfig().GrpcAddress)
	}

	return c.repo.userRepo
}

func (c *serviceProvider) GetUserValidator(_ context.Context) usecase.UserValidator {
	if c.validator.userValidator == nil {
		c.validator.userValidator = validator.NewUserValidator()
	}

	return c.validator.userValidator
}

func (c *serviceProvider) GetCreateUserUseCase(ctx context.Context) usecase.CreateUserUseCase {
	if c.useCase.createUserUseCase == nil {
		c.useCase.createUserUseCase = usecase.NewCreateUserUseCase(c.GetUserValidator(ctx), c.GetUserRepo(ctx), c.GetLogger(ctx))
	}

	return c.useCase.createUserUseCase
}

func (c *serviceProvider) GetListUserUseCase(ctx context.Context) usecase.ListUserUseCase {
	if c.useCase.userListUseCase == nil {
		c.useCase.userListUseCase = usecase.NewListUserUseCase(c.GetUserRepo(ctx), c.GetLogger(ctx))
	}

	return c.useCase.userListUseCase
}
