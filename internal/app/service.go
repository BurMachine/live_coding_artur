package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"awesomeProject/internal/app/counter"
	"awesomeProject/internal/config"
	"awesomeProject/internal/transport"
	"awesomeProject/internal/transport/handlers"
)

func Run(ctx context.Context, cfg *config.Config) error {
	ctx, cancel := context.WithCancel(ctx)

	s := transport.New(cfg)
	c := counter.New()
	h := handlers.New(cfg, c)

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
