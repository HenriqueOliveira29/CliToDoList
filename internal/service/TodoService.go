package service

import (
	"strings"
	"time"
	"context"
	"errors"
	"github.com/HenriqueOliveira29/cliToDoList/internal/repository"
	"github.com/HenriqueOliveira29/cliToDoList/internal/model"
)

var ErrInvalidTitle = errors.New("title is required")

type TodoService struct {
	repo repository.ToDoRepository
}

func NewToDoService(repo repository.ToDoRepository) *TodoService{
	return &TodoService{repo: repo}
}

func(s *TodoService) GetAll(ctx context.Context) ([]model.Task, error){
	return s.repo.GetAll(ctx)
}

func(s *TodoService) GetById(ctx context.Context, id int) (model.Task, error){
	return s.repo.GetById(ctx, id)
}

func(s *TodoService) Add(ctx context.Context,title string) (int, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return 0, ErrInvalidTitle
	}
	return s.repo.Add(ctx, model.Task{
		Title : title,
		CreatedDate: time.Now(),
	})
}

func(s *TodoService) UpdateTitle(ctx context.Context, id int, title string) error {
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	task.Title = title
	return s.repo.Update(ctx, task)
}

func(s *TodoService) MarkAsDone(ctx context.Context, id int) error {
	task, err := s.repo.GetById(ctx, id)
	if err != nil {
		return err
	}
	task.Done = true
	return s.repo.Update(ctx, task)
}

func(s *TodoService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
