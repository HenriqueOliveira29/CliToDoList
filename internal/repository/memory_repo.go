package repository

import (
	"fmt"
	"github.com/HenriqueOliveira29/cliToDoList/internal/model"
)

type MemoryRepo struct {
	Tasks []model.Task 
}

func NewMemoryRepo() *MemoryRepo {
	return &MemoryRepo{}
}

func (r *MemoryRepo) GetAll() ([]model.Task, error) {
	return r.Tasks, nil
} 	

func (r *MemoryRepo) GetById(id int) (model.Task, error) {
	for _, element := range r.Tasks {
		if element.Id == id {
			return element, nil
		} 
	}
	return model.Task{}, fmt.Errorf("We dont have a task with the Id %d", id)
}

func (r *MemoryRepo) Add(task model.Task) error {
	r.Tasks = append(r.Tasks, task)
	return nil
}

func (r *MemoryRepo) Update(task model.Task) error {
	for _, element := range r.Tasks {
		if element.Id == task.Id {
			element = task
			return nil
		}
	}
	return fmt.Errorf("We dont have this task on your list of tasks")
}

func (r *MemoryRepo) Delete(id int) error {
	for i, t := range r.Tasks {
		if t.Id == id {
			r.Tasks = append(r.Tasks[:i], r.Tasks[i+1])
			return nil
		}
	}
	return fmt.Errorf("We dont have this task on your list of tasks")
}
