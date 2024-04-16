package main

import (
	"context"
	"diploma-1/internal/config"
	"fmt"
)

func main() {
	ctx := context.Background()
	conf, err := config.New(ctx)
	if err != nil {
		fmt.Printf("unable to collect config: %v", err)
	}
	fmt.Printf("applied args: %s", conf)

}
