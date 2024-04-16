package main

import (
	"context"
	"diploma-1/internal/config"
	"diploma-1/internal/logger"
	"diploma-1/internal/storage"
	"fmt"
)

// TODO: добавить middleware с logger.With() для добавления трейса запроса

func main() {
	ctx := context.Background()
	conf, err := config.New(ctx)
	if err != nil {
		fmt.Printf("unable to collect config: %v", err)
	}
	fmt.Printf("applied args: %s\n", conf)
	logger.New(conf)
	if err := storage.New(ctx); err != nil {
		logger.Fatalf(ctx, "unable to init storage: %v", err)
	}
}
