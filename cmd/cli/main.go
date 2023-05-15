package main

import (
	"context"
	"log"

	"github.com/a13hander/chat-server/internal/app"
)

func main() {
	ctx := context.Background()
	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	err = a.Run(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}
