package app

import (
	"context"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/a13hander/chat-server/internal/config"
	"github.com/a13hander/chat-server/internal/interceptor"
	desc "github.com/a13hander/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type App struct {
	serviceProvider *serviceProvider
	config          *config.Config
	grpcServer      *grpc.Server
	swaggerServer   *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.closeAll()
		closer.wait()
	}()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := a.runGrpcServer(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGrpcServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	a.config = config.GetConfig()
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGrpcServer(ctx context.Context) error {
	accessInter := interceptor.NewAuthInterceptor(a.serviceProvider.GetAccessChecker(ctx))

	a.grpcServer = grpc.NewServer(grpc.UnaryInterceptor(accessInter.Unary()))

	desc.RegisterChatV1Server(a.grpcServer, a.serviceProvider.GetChatV1ServerImpl(ctx))

	reflection.Register(a.grpcServer)

	return nil
}

func (a *App) runGrpcServer(_ context.Context) error {
	log.Printf("Grpc server starting on %s\n", a.config.GrpcPort)

	listener, err := net.Listen("tcp", a.config.GrpcPort)
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(listener)
	if err != nil {
		return err
	}

	return nil
}
