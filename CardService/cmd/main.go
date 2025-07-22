package main

import (
	"CardService/internal/application"
	"context"
	"golang.org/x/sync/errgroup"
)

func main() {
	errG, ctx := errgroup.WithContext(context.Background())
	ctx, shutdownHook := context.WithCancel(ctx)

	app := application.NewApp(ctx)

	app.Start(errG)

	app.AwaitTermination(shutdownHook)

}
