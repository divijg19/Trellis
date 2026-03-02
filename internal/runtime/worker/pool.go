package worker

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/divijg19/Trellis/internal/domain"
	"github.com/divijg19/Trellis/internal/runtime/handlers"
)

type TaskRepository interface {
	GetByID(id string) (*domain.Task, error)
	Update(task *domain.Task) error
}

type Pool struct {
	concurrency int
	queue       <-chan string
	repository  TaskRepository
	registry    *handlers.Registry
	logger      *log.Logger

	workers sync.WaitGroup
}

func NewPool(
	concurrency int,
	queue <-chan string,
	repository TaskRepository,
	registry *handlers.Registry,
	logger *log.Logger,
) *Pool {
	return &Pool{
		concurrency: concurrency,
		queue:       queue,
		repository:  repository,
		registry:    registry,
		logger:      logger,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := 0; i < p.concurrency; i++ {
		p.workers.Add(1)
		go p.runWorker(ctx, i+1)
	}
}

func (p *Pool) Wait() {
	p.workers.Wait()
}

func (p *Pool) runWorker(ctx context.Context, workerID int) {
	defer p.workers.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case taskID := <-p.queue:
			p.processTask(ctx, workerID, taskID)
		}
	}
}

func (p *Pool) processTask(ctx context.Context, workerID int, taskID string) {
	task, err := p.repository.GetByID(taskID)
	if err != nil {
		p.logger.Printf("worker=%d task=%s load failed: %v", workerID, taskID, err)
		return
	}

	if err := task.TransitionTo(domain.TaskStatusRunning, time.Now().UTC()); err != nil {
		p.logger.Printf("worker=%d task=%s transition queued->running failed: %v", workerID, taskID, err)
		return
	}
	if err := p.repository.Update(task); err != nil {
		p.logger.Printf("worker=%d task=%s update running failed: %v", workerID, taskID, err)
		return
	}

	handler, ok := p.registry.Get(task.Type)
	if !ok {
		p.markFailed(task, workerID, "handler not found")
		return
	}

	err = handler(ctx, task.Payload)
	if err != nil {
		p.markFailed(task, workerID, err.Error())
		return
	}

	if err := task.TransitionTo(domain.TaskStatusCompleted, time.Now().UTC()); err != nil {
		p.logger.Printf("worker=%d task=%s transition running->completed failed: %v", workerID, taskID, err)
		return
	}

	if err := p.repository.Update(task); err != nil {
		p.logger.Printf("worker=%d task=%s update completed failed: %v", workerID, taskID, err)
		return
	}

	p.logger.Printf("worker=%d task=%s completed", workerID, taskID)
}

func (p *Pool) markFailed(task *domain.Task, workerID int, reason string) {
	if err := task.TransitionTo(domain.TaskStatusFailed, time.Now().UTC()); err != nil {
		p.logger.Printf("worker=%d task=%s transition running->failed failed: %v", workerID, task.ID, err)
		return
	}

	if err := p.repository.Update(task); err != nil {
		p.logger.Printf("worker=%d task=%s update failed status failed: %v", workerID, task.ID, err)
		return
	}

	p.logger.Printf("worker=%d task=%s failed: %s", workerID, task.ID, reason)
}
