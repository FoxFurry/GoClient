package main

import (
	"context"
	"github.com/foxfurry/go_client/application"
	"github.com/foxfurry/go_client/internal/infrastructure/config"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	config.LoadConfig()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	app := application.CreateApp(ctx)
	go app.Start()

	<-sigChan

	app.Shutdown()
	cancel()
}
