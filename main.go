package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/juho05/log"
	"github.com/juho05/stine-ical-formatter/web"
)

func run() error {
	port := 8080
	if os.Getenv("PORT") != "" {
		var err error
		port, err = strconv.Atoi(os.Getenv("PORT"))
		if err != nil {
			return fmt.Errorf("invalid PORT value: %w", err)
		}
	}
	server, err := web.NewServer(fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("new server: %w", err)
	}
	closed := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		timeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
		err := server.Shutdown(timeout)
		if err != nil {
			log.Error("shutdown server: %w", err)
		}
		cancelTimeout()
		close(closed)
	}()

	err = server.Listen()
	if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	if err == nil {
		<-closed
	}
	return err
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
