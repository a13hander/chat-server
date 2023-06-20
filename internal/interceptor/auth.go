package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/a13hander/chat-server/internal/domain/usecase"
)

type AuthInterceptor struct {
	checkAccessUseCase usecase.CheckAccessUseCase
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if ok {
			ctx = metadata.NewOutgoingContext(ctx, md)
		}

		check, err := a.checkAccessUseCase.Run(ctx, info.FullMethod)
		if err != nil || !check {
			return nil, fmt.Errorf("access deniend: %w", err)
		}

		response, err := handler(ctx, req)
		if err != nil {
			return nil, err
		}

		return response, nil
	}
}
