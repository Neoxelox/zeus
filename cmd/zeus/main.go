package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/neoxelox/zeus/internal/server"
)

func main() {
	instance := echo.New()
	zeus := server.New(instance)
	go zeus.Startup()

	// Graceful shutdown.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(
		context.Background(), time.Duration(zeus.Configuration.App.GracefulTimeout)*time.Second)
	defer cancel()
	if err := zeus.Shutdown(ctx); err != nil {
		log.Fatal(err) // nolint
	}
}
