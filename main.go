package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/salmanahmed404/go-keyre/internal/server"
)

func main() {
	service := ":1200"
	sig := make(chan os.Signal)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	srv := server.NewServer(service)

	select {
	case <-sig:
		srv.Stop()
	}
}
