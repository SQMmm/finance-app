package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/sqmmm/finance-app/internal/config"
	"github.com/sqmmm/finance-app/internal/container"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Configuration read error: %v", err)
	}

	if err = container.Build(cfg); err != nil {
		log.Fatalf("Initialization error: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	server := container.GetServer()
	wg := &sync.WaitGroup{}

	go func() {
		err = server.Serve(ctx, wg)
		if err != nil {
			log.Fatalf("failed to serve: %s", err)
		}
	}()

	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	<-gracefulStop
	cancel()
	wg.Wait()
}
