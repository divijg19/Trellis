package runtime

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/divijg19/Trellis/internal/domain"
	"github.com/divijg19/Trellis/internal/queue"
	"github.com/divijg19/Trellis/internal/runtime/handlers"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	GetByID(id string) (*domain.Task, error)
	List() ([]*domain.Task, error)
	Update(task *domain.Task) error
}

var (
	ErrInvalidTaskType = errors.New("invalid task type")
)

type TaskService struct {
	repository TaskRepository
	queue      *queue.TaskQueue
	registry   *handlers.Registry
}

func NewTaskService(repository TaskRepository, queue *queue.TaskQueue, registry *handlers.Registry) *TaskService {
	return &TaskService{repository: repository, queue: queue, registry: registry}
}

func (s *TaskService) CreateTask(taskType string, payload []byte) (*domain.Task, error) {
	if _, ok := s.registry.Get(taskType); !ok {
		return nil, ErrInvalidTaskType
	}

	now := time.Now().UTC()
	task := &domain.Task{
		ID:        newTaskID(),
		Type:      taskType,
		Payload:   append([]byte(nil), payload...),
		Status:    domain.TaskStatusQueued,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repository.Create(task); err != nil {
		return nil, err
	}

	s.queue.Enqueue(task.ID)
	return task, nil
}

func (s *TaskService) GetTask(id string) (*domain.Task, error) {
	return s.repository.GetByID(id)
}

func (s *TaskService) ListTasks() ([]*domain.Task, error) {
	return s.repository.List()
}

func newTaskID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return hex.EncodeToString([]byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}

	return hex.EncodeToString(b)
}
