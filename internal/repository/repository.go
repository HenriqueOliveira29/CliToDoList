package repository

import (
	"context"
	"github.com/HenriqueOliveira29/cliToDoList/internal/model"
)

type ToDoRepository interface {
	GetAll(ctx context.Context) ([]model.Task, error)
	Add(ctx context.Context, task model.Task) (int, error)
	GetById(ctx context.Context, id int) (model.Task, error)
	Update(ctx context.Context, task model.Task) error
	Delete(ctx context.Context, id int) error
}
