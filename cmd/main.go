package main

import (
	"context"
	"log"

	"github.com/BelyaevEI/microservices_auth/internal/app"
)

func main() {

	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	go func() {
		err = a.RunConsumerForCreateUser(ctx)
		if err != nil {
			log.Fatalf("failed to consume saver: %s", err.Error())
		}
	}()

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
