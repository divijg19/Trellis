package handlers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type HandlerFunc func(ctx context.Context, payload []byte) error

type Registry struct {
	mu       sync.RWMutex
	handlers map[string]HandlerFunc
}

func NewRegistry() *Registry {
	return &Registry{handlers: make(map[string]HandlerFunc)}
}

func (r *Registry) Register(taskType string, handler HandlerFunc) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.handlers[taskType] = handler
}

func (r *Registry) Get(taskType string) (HandlerFunc, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	handler, ok := r.handlers[taskType]
	return handler, ok
}

func RegisterDefaultHandlers(registry *Registry, logger *log.Logger) {
	registry.Register("echo", func(ctx context.Context, payload []byte) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			logger.Printf("echo handler payload=%s", string(payload))
			return nil
		}
	})

	registry.Register("sleep", func(ctx context.Context, payload []byte) error {
		raw := strings.TrimSpace(string(payload))
		seconds, err := strconv.Atoi(raw)
		if err != nil {
			return fmt.Errorf("invalid sleep payload %q: %w", raw, err)
		}
		if seconds < 0 {
			return fmt.Errorf("sleep duration must be non-negative")
		}

		timer := time.NewTimer(time.Duration(seconds) * time.Second)
		defer timer.Stop()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-timer.C:
			return nil
		}
	})
}
