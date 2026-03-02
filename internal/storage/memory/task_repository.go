package memory

import (
	"sort"
	"sync"

	"github.com/divijg19/Trellis/internal/domain"
)

type TaskRepository struct {
	mu    sync.RWMutex
	tasks map[string]*domain.Task
}

func NewTaskRepository() *TaskRepository {
	return &TaskRepository{tasks: make(map[string]*domain.Task)}
}

func (r *TaskRepository) Create(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tasks[task.ID] = cloneTask(task)
	return nil
}

func (r *TaskRepository) GetByID(id string) (*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	task, ok := r.tasks[id]
	if !ok {
		return nil, domain.ErrTaskNotFound
	}

	return cloneTask(task), nil
}

func (r *TaskRepository) List() ([]*domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]*domain.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, cloneTask(task))
	}

	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].CreatedAt.Before(tasks[j].CreatedAt)
	})

	return tasks, nil
}

func (r *TaskRepository) Update(task *domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.tasks[task.ID]; !ok {
		return domain.ErrTaskNotFound
	}

	r.tasks[task.ID] = cloneTask(task)
	return nil
}

func cloneTask(task *domain.Task) *domain.Task {
	if task == nil {
		return nil
	}

	clone := *task
	if task.Payload != nil {
		clone.Payload = append([]byte(nil), task.Payload...)
	}

	return &clone
}
