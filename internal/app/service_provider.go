package app

import (
	"context"
	"log"

	accessV1 "github.com/a13hander/auth-service-api/pkg/access_v1"
	authV1 "github.com/a13hander/auth-service-api/pkg/auth_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/a13hander/chat-server/internal/config"
	"github.com/a13hander/chat-server/internal/domain/usecase"
	"github.com/a13hander/chat-server/internal/domain/util"
	"github.com/a13hander/chat-server/internal/domain/validator"
	"github.com/a13hander/chat-server/internal/service/auth"
	"github.com/a13hander/chat-server/internal/service/logger"
)

type serviceProvider struct {
	logger       util.Logger
	grpcClient   *grpc.ClientConn
	accessClient accessV1.AccessV1Client
	authClient   authV1.AuthV1Client

	repo struct {
		userRepo usecase.UserRepo
	}

	validator struct {
		userValidator usecase.UserValidator
	}

	useCase struct {
		createUserUseCase usecase.CreateUserUseCase
		userListUseCase   usecase.ListUserUseCase

		checkAccessUseCase usecase.CheckAccessUseCase
	}

	service struct {
		accessChecker usecase.AccessChecker
	}
}

func NewServiceProvider() *serviceProvider {
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
		c.repo.userRepo = auth.NewAuthClient(c.GetAuthClient())
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

func (c *serviceProvider) GetGrpcClient() *grpc.ClientConn {
	if c.grpcClient == nil {
		conn, err := grpc.Dial(config.GetConfig().GrpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalln(err)
		}

		c.grpcClient = conn
	}

	return c.grpcClient
}

func (c *serviceProvider) GetAuthClient() authV1.AuthV1Client {
	if c.authClient == nil {
		c.authClient = authV1.NewAuthV1Client(c.GetGrpcClient())
	}

	return c.authClient
}

func (c *serviceProvider) GetAccessClient() accessV1.AccessV1Client {
	if c.accessClient == nil {
		c.accessClient = accessV1.NewAccessV1Client(c.GetGrpcClient())
	}

	return c.accessClient
}

func (c *serviceProvider) GetAccessChecker() usecase.AccessChecker {
	if c.service.accessChecker == nil {
		c.service.accessChecker = auth.NewAccessChecker(c.GetAccessClient())
	}

	return c.service.accessChecker
}

func (c *serviceProvider) GetCheckAccessUseCase() usecase.CheckAccessUseCase {
	if c.useCase.checkAccessUseCase == nil {
		c.useCase.checkAccessUseCase = usecase.NewCheckAccessUseCase(c.GetAccessChecker())
	}

	return c.useCase.checkAccessUseCase
}
