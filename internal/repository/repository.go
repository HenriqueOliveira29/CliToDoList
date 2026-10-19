package repository

import (
	"github.com/HenriqueOliveira29/cliToDoList/internal/model"
)

type ToDoRepository interface {
	GetAll() ([]model.Task, error)
	Add(task model.Task) error
	GetById(id int) (model.Task, error)
	Update(task model.Task) error
	Delete(id int)
}
