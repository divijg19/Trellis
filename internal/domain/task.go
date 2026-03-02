package domain

import (
	"errors"
	"fmt"
	"time"
)

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusQueued    TaskStatus = "queued"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
)

var ErrInvalidStatusTransition = errors.New("invalid task status transition")
var ErrTaskNotFound = errors.New("task not found")

type Task struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Payload   []byte     `json:"payload"`
	Status    TaskStatus `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func (t *Task) CanTransitionTo(next TaskStatus) bool {
	switch t.Status {
	case TaskStatusPending:
		return next == TaskStatusQueued
	case TaskStatusQueued:
		return next == TaskStatusRunning
	case TaskStatusRunning:
		return next == TaskStatusCompleted || next == TaskStatusFailed
	case TaskStatusCompleted, TaskStatusFailed:
		return false
	default:
		return false
	}
}

func (t *Task) TransitionTo(next TaskStatus, updatedAt time.Time) error {
	if !t.CanTransitionTo(next) {
		return fmt.Errorf("%w: %s -> %s", ErrInvalidStatusTransition, t.Status, next)
	}

	t.Status = next
	t.UpdatedAt = updatedAt
	return nil
}
