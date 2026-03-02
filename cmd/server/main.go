package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/divijg19/Trellis/internal/api/httpapi"
	"github.com/divijg19/Trellis/internal/queue"
	"github.com/divijg19/Trellis/internal/runtime"
	"github.com/divijg19/Trellis/internal/runtime/handlers"
	"github.com/divijg19/Trellis/internal/runtime/worker"
	"github.com/divijg19/Trellis/internal/storage/memory"
)

const (
	addr              = ":8080"
	workerConcurrency = 3
	queueBuffer       = 256
)

func main() {
	logger := log.New(os.Stdout, "trellis ", log.LstdFlags)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	repository := memory.NewTaskRepository()
	taskQueue := queue.NewTaskQueue(queueBuffer)

	registry := handlers.NewRegistry()
	handlers.RegisterDefaultHandlers(registry, logger)

	service := runtime.NewTaskService(repository, taskQueue, registry)
	apiServer := httpapi.NewServer(service)

	pool := worker.NewPool(workerConcurrency, taskQueue.Consume(), repository, registry, logger)
	pool.Start(ctx)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: apiServer.Handler(),
	}

	go func() {
		<-ctx.Done()
		if err := httpServer.Shutdown(context.Background()); err != nil {
			logger.Printf("http shutdown error: %v", err)
		}
	}()

	logger.Printf("server listening on %s", addr)
	if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalf("http server failed: %v", err)
	}

	pool.Wait()
}
