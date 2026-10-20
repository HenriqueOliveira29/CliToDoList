package api

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/HenriqueOliveira29/cliToDoList/internal/service"
)

type Server struct {
	svc *service.TodoService
}

func NewServer(svc *service.TodoService) *Server {
	return &Server{svc: svc}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /tasks", s.listTasks)
	mux.HandleFunc("POST /tasks/add", s.createTask)
	return mux
}

// Start blocks until the server fails or ctx is cancelled, then shuts down gracefully.
func (s *Server) Start(ctx context.Context, addr string) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           s.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}