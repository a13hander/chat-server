package app

import (
	"context"

	"github.com/a13hander/chat-server/internal/app/cli"
	"github.com/a13hander/chat-server/internal/config"
)

type App struct {
	serviceProvider *serviceProvider
	config          *config.Config
	appCli          cli.AppCli
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
	return a.appCli.Run(ctx)
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initAppCli,
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
	a.serviceProvider = NewServiceProvider()
	return nil
}

func (a *App) initAppCli(ctx context.Context) error {
	if a.appCli == nil {
		a.appCli = cli.NewAppCli(a.serviceProvider.GetCreateUserUseCase(ctx), a.serviceProvider.GetListUserUseCase(ctx))
	}

	return nil
}

func (a *App) RunCli(ctx context.Context) error {
	return a.appCli.Run(ctx)
}
