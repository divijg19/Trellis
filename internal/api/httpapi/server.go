package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/divijg19/Trellis/internal/domain"
	"github.com/divijg19/Trellis/internal/runtime"
)

type Server struct {
	service *runtime.TaskService
	mux     *http.ServeMux
}

type createTaskRequest struct {
	Type    string `json:"type"`
	Payload string `json:"payload"`
}

type taskResponse struct {
	ID        string            `json:"id"`
	Type      string            `json:"type"`
	Payload   string            `json:"payload"`
	Status    domain.TaskStatus `json:"status"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

func NewServer(service *runtime.TaskService) *Server {
	s := &Server{service: service, mux: http.NewServeMux()}
	s.routes()
	return s
}

func (s *Server) Handler() http.Handler {
	return s.mux
}

func (s *Server) routes() {
	s.mux.HandleFunc("/tasks", s.handleTasks)
	s.mux.HandleFunc("/tasks/", s.handleTaskByID)
}

func (s *Server) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.createTask(w, r)
	case http.MethodGet:
		s.listTasks(w, r)
	default:
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "invalid json body")
		return
	}

	req.Type = strings.TrimSpace(req.Type)
	if req.Type == "" {
		s.respondError(w, http.StatusBadRequest, "type is required")
		return
	}

	task, err := s.service.CreateTask(req.Type, []byte(req.Payload))
	if err != nil {
		if errors.Is(err, runtime.ErrInvalidTaskType) {
			s.respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		s.respondError(w, http.StatusInternalServerError, "failed to create task")
		return
	}

	s.respondJSON(w, http.StatusCreated, toTaskResponse(task))
}

func (s *Server) listTasks(w http.ResponseWriter, _ *http.Request) {
	tasks, err := s.service.ListTasks()
	if err != nil {
		s.respondError(w, http.StatusInternalServerError, "failed to list tasks")
		return
	}

	res := make([]taskResponse, 0, len(tasks))
	for _, task := range tasks {
		res = append(res, toTaskResponse(task))
	}

	s.respondJSON(w, http.StatusOK, res)
}

func (s *Server) handleTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		s.respondError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/tasks/")
	id = strings.TrimSpace(id)
	if id == "" {
		s.respondError(w, http.StatusNotFound, "task id required")
		return
	}

	task, err := s.service.GetTask(id)
	if err != nil {
		if errors.Is(err, domain.ErrTaskNotFound) {
			s.respondError(w, http.StatusNotFound, err.Error())
			return
		}

		s.respondError(w, http.StatusInternalServerError, "failed to get task")
		return
	}

	s.respondJSON(w, http.StatusOK, toTaskResponse(task))
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	s.respondJSON(w, status, map[string]string{"error": message})
}

func toTaskResponse(task *domain.Task) taskResponse {
	return taskResponse{
		ID:        task.ID,
		Type:      task.Type,
		Payload:   string(task.Payload),
		Status:    task.Status,
		CreatedAt: task.CreatedAt.UTC().Format(timeLayout),
		UpdatedAt: task.UpdatedAt.UTC().Format(timeLayout),
	}
}

const timeLayout = "2006-01-02T15:04:05.999999999Z07:00"
