package main

import (
	"context"
	"log"

	"github.com/sfomuseum/go-sfomuseum-github/app/ratelimit/monitor"
)

func main() {

	ctx := context.Background()
	logger := log.Default()

	err := monitor.Run(ctx, logger)

	if err != nil {
		logger.Fatalf("Failed to run monitor, %v", err)
	}
}
