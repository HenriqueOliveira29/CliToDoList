package api

import (
	"net/http"
	"github.com/HenriqueOliveira29/cliToDoList/internal/service"
)

type Server struct {
	svc *service.TodoService
}

func NewServer(svc *service.TodoService) *Server {
	return &Server{svc: svc}
} 

func (s *Server) Start(port string) error {
	mux:=http.NewServeMux()
	mux.HandleFunc("GET /tasks", s.listTasks)
	mux.HandleFunc("POST /tasks/add", s.createTask)

	return http.ListenAndServe(":"+port, mux)
}


