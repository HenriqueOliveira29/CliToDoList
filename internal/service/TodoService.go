package service

import (
	"strings"
	"time"
	"github.com/HenriqueOliveira29/cliToDoList/internal/repository"
	"github.com/HenriqueOliveira29/cliToDoList/internal/model"
)

type TodoService struct {
	repo *repository.MemoryRepo
}

func NewToDoService(repo *repository.MemoryRepo) *TodoService{
	return &TodoService{repo: repo}
}

func(s *TodoService) GetAll() ([]model.Task, error){
	return s.repo.GetAll()
}

func(s *TodoService) GetById(id int) (model.Task, error){
	return s.repo.GetById(id)
}

func(s *TodoService) Add(title string) error {
	title = strings.TrimSpace(title)
	s.repo.Add(model.Task{
		Title : title,
		CreatedDate: time.Now(),
	})
	return nil
}

func(s *TodoService) UpdateTitle(id int, title string) error {
	task, err := s.repo.GetById(id)
	if err != nil {
		return err
	}
	task.Title = title
	return s.repo.Update(task)
}

func(s *TodoService) MarkAsDone(id int) error {
	task, err := s.repo.GetById(id)
	if err != nil {
		return err
	}
	task.Done = true
	return s.repo.Update(task)
}

func(s *TodoService) Delete(id int) error {
	return s.repo.Delete(id)
}
