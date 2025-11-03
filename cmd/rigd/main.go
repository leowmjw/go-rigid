package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/leowmjw/go-rigid/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	mux := http.NewServeMux()
	server.RegisterHandlers(mux, logger)

	srv := &http.Server{ Addr: ":8080", Handler: mux }

	go func() {
		logger.Info("http listen", "addr", srv.Addr)
		_ = srv.ListenAndServe()
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()
	_ = srv.Shutdown(context.Background())
}
