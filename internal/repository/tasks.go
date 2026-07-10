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
		"INSERT INTO tasks (title, status,group_id) VALUES (?,?,?)",
		task.Title,
		task.Status,
		task.GroupId)
	if err != nil {
		return err
	}
	task.ID, _ = result.LastInsertId()
	return nil
}

func (r *TaskRepository) List(ctx context.Context, groupID int64, orderBy string) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,title,status,created_at,group_id
		 FROM tasks
		WHERE group_id = ?
			ORDER BY
				CASE WHEN ? = 'done' THEN status END DESC,
				CASE WHEN ? = 'created_at' THEN created_at END ASC,
				id ASC
		`, groupID, orderBy, orderBy)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.CreatedAt, &t.GroupId); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *TaskRepository) Complete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET status = ? WHERE id = ?", "DONE", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET status = ? WHERE id = ?", status, id)
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

func (r *TaskRepository) DeleteCompleted(ctx context.Context, groupID int64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM tasks WHERE status = ? and group_id = ?", "DONE", groupID)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) Truncate(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM tasks")
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskRepository) GetById(ctx context.Context, id int64) (*models.Task, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id,title,status,created_at,group_id FROM tasks WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var task models.Task

	if err := row.Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt, &task.GroupId); err != nil {
		return nil, err
	}
	return &task, nil
}
