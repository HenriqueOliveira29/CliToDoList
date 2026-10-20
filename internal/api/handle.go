package api

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/HenriqueOliveira29/cliToDoList/internal/service"
)

type createTaskRequest struct {
	Title string `json:"title"`
}

type createTaskResponse struct {
	ID int `json:"id"`
}

func (s *Server) listTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := s.svc.GetAll(r.Context())
	if err != nil {
		slog.Error("list tasks", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, tasks)
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1 MiB

	var req createTaskRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	id, err := s.svc.Add(r.Context(), req.Title)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTitle) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		slog.Error("create task", "err", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusCreated, createTaskResponse{ID: id})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("write response", "err", err)
	}
}