package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/HenriqueOliveira29/cliToDoList/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("task not found")

type PostgresRepo struct {
	db *pgxpool.Pool
}

func NewPostgresRepo(db *pgxpool.Pool) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) GetAll(ctx context.Context) ([]model.Task, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, title, done, created_at AS "CreatedDate" FROM todos ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("query todos: %w", err)
	}

	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Task])
	if err != nil {
		return nil, fmt.Errorf("collect todos: %w", err)
	}
	return tasks, nil
}

func (r *PostgresRepo) GetById(ctx context.Context, id int) (model.Task, error) {
	var task model.Task
	err := r.db.QueryRow(ctx,
		`SELECT id, title, done, created_at AS "CreatedDate" FROM todos WHERE id = $1`, id).Scan(&task.Id, &task.Title, &task.Done, &task.CreatedDate)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, ErrNotFound
		}
		return model.Task{}, fmt.Errorf("query todo: %w", err)
	}
	return task, nil
}

func (r *PostgresRepo) Add(ctx context.Context, task model.Task) (int, error) {
	var id int
	err := r.db.QueryRow(ctx,
		`INSERT INTO todos (title, done) VALUES ($1, $2) RETURNING id`,
		task.Title, task.Done,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("insert todo: %w", err)
	}
	return id, nil
}

func (r *PostgresRepo) Update(ctx context.Context, task model.Task) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE todos SET title = $1, done = $2 WHERE id = $3`,
		task.Title, task.Done, task.Id,
	)
	if err != nil {
		return fmt.Errorf("update todo: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, id int) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM todos WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete todo: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}