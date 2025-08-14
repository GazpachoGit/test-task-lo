package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"test-task-lo/internal/service/asynclog"
)

func main() {
	//TODO: add config file, add log level

	//start async logger
	log := asynclog.NewAsyncLog(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})), 3)
	asynclog.StartLogger(log)

	// log.Info("test 1")
	// log.Error("test 2")
	// log.Warn("test 3")
	// log.Debug("test 4")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("Received shutdown signal", slog.String("signal", sign.String()))
	asynclog.StopLogger(log)
}
