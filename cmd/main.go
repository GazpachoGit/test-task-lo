package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	gettermulty "test-task-lo/internal/http-server/handlers/task/getter-multy"
	getterone "test-task-lo/internal/http-server/handlers/task/getter-one"
	"test-task-lo/internal/http-server/handlers/task/setter"
	sl "test-task-lo/internal/lib/log"
	"test-task-lo/internal/service/asynclog"
	storage "test-task-lo/internal/storage/inmemory-storage"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	//TODO: add a config file, for a log level

	//handler can be replaced with JSON, DB,...
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})

	log := asynclog.NewAsyncLog(handler, 3)
	wg := &sync.WaitGroup{}
	//start async logger
	asynclog.StartLogger(log, wg)
	log.Debug("Logger ready")

	storage := storage.New()

	router := chi.NewRouter()
	router.Use(middleware.URLFormat)

	router.Post("/tasks", setter.NewSet(log, storage))
	router.Get("/tasks", gettermulty.New(log, storage))
	router.Get("/tasks/{id}", getterone.New(log, storage))

	srv := startHTTPServer(router, wg, log)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	sign := <-stop

	//graceful stop
	log.Info("Received shutdown signal", slog.String("signal", sign.String()))
	stopHTTPServer(srv, log)
	log.Info("Stopping logger")
	asynclog.StopLogger(log)
	wg.Wait()
	fmt.Println("Graceful shutdown complete")
}

func startHTTPServer(router *chi.Mux, wg *sync.WaitGroup, log *asynclog.AsyncLog) *http.Server {
	//TODO: add a config file
	srv := &http.Server{
		Addr:         "localhost:9000",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info("Starting server on :9000")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("can't start http server", sl.Err(err))
		}
	}()
	return srv
}

func stopHTTPServer(srv *http.Server, log *asynclog.AsyncLog) {
	log.Info("Stopping HTTP Server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("HTTP server shutdown error", sl.Err(err))
	}
	log.Info("HTTP Server stoped")
}
