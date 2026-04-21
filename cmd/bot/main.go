package main

import (
	"context"
	"log"
	"notification-service/internal/app"
)

func main() {
	ctx := context.Background()
	app, err := app.NewApp(ctx)
	if err != nil {
		panic(err)
	}
	if err := app.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
