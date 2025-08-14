package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"test-task-lo/internal/http-server/handlers/task/setter"
	"test-task-lo/internal/service/asynclog"
	storage "test-task-lo/internal/storage/inmemory-storage"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	//TODO: add config file, add log level

	//start async logger
	//hadler can be replaced with JSON, DB,...
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	log := asynclog.NewAsyncLog(handler, 3)
	asynclog.StartLogger(log)
	log.Debug("Logger ready")

	// log.Info("test 1")
	// log.Error("test 2")
	// log.Warn("test 3")
	// log.Debug("test 4")

	storage := storage.New()

	router := chi.NewRouter()
	router.Use(middleware.URLFormat)

	router.Post("/tasks", setter.NewSet(log, storage))
	router.Get("/tasks/{id}", getter.New(log, storage))

	//stop logger
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop
	log.Info("Received shutdown signal", slog.String("signal", sign.String()))
	asynclog.StopLogger(log)
}
