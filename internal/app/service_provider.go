package app

import (
	"context"
	"log"

	accessV1 "github.com/a13hander/auth-service-api/pkg/access_v1"
	"github.com/a13hander/chat-server/internal/access"
	chatV1 "github.com/a13hander/chat-server/internal/app/chat_v1"
	"github.com/a13hander/chat-server/internal/config"
	"github.com/a13hander/chat-server/internal/interceptor"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type serviceProvider struct {
	accessClient     accessV1.AccessV1Client
	chatV1ServerImpl *chatV1.Implementation
	accessChecker    interceptor.AccessChecker
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (c *serviceProvider) GetAccessClient() accessV1.AccessV1Client {
	if c.accessClient == nil {
		conn, err := grpc.Dial(config.GetConfig().AccessAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalln(err)
		}

		c.accessClient = accessV1.NewAccessV1Client(conn)
	}

	return c.accessClient
}

func (c *serviceProvider) GetChatV1ServerImpl(ctx context.Context) *chatV1.Implementation {
	return chatV1.NewImplementation()
}

func (c *serviceProvider) GetAccessChecker(ctx context.Context) interceptor.AccessChecker {
	if c.accessChecker == nil {
		c.accessChecker = access.NewAccessChecker(c.GetAccessClient())
	}

	return c.accessChecker
}
