package repository

import (
	"context"
	"database/sql"

	"github.com/fazil-syed/todoer/internal/models"
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
		"SELECT id,title,done,created_at FROM tasks")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.ID, &t.Title, &t.Done, &t.CreatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepository) Complete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET done = ? WHERE id = ?", true, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) DeleteCompleted(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM tasks WHERE done = ?", true)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) GetById(ctx context.Context, id int64) (*models.Task, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id,title,done,created_at FROM tasks WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var task models.Task

	if err := row.Scan(&task.ID, &task.Title, &task.Done, &task.CreatedAt); err != nil {
		return nil, err
	}
	return &task, nil
}
