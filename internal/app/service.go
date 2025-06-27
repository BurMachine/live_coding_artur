package app

import (
	"awesomeProject/internal/config"
	"awesomeProject/internal/transport"
	"awesomeProject/internal/transport/handlers"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(ctx)

	s := transport.New(cfg)
	h := handlers.New(cfg)

	go s.Start(ctx, h)

	gracefulShutDown(cancel)
	time.Sleep(2 * time.Second)
	return nil
}

func gracefulShutDown(cancel context.CancelFunc) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)

	<-ch
	cancel()

}
