package repository

import (
	"context"
	"database/sql"

	"github.com/fazil-syed/todoer/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}
func (r *TaskRepository) Create(ctx context.Context, task *models.Task) error {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO tasks (title, done) VALUES (?,?)",
		task.Title,
		task.Done)
	if err != nil {
		return err
	}
	task.ID, _ = result.LastInsertId()
	return nil
}

func (r *TaskRepository) List(ctx context.Context) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT * FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
