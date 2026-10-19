package api

import (
    "encoding/json"
    "net/http"
)

type createTaskRequest struct {
    Title string `json:"title"`
}

func (s *Server) listTasks(w http.ResponseWriter, r *http.Request) {
    tasks, err := s.svc.GetAll()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}

func (s *Server) createTask(w http.ResponseWriter, r *http.Request) {
    var req createTaskRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "invalid JSON body", http.StatusBadRequest)
        return
    }
    if err := s.svc.Add(req.Title); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
}
