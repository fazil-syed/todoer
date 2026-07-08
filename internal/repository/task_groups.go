package repository

import (
	"context"
	"database/sql"

	"github.com/fazil-syed/todoer/internal/models"
)

type TaskGroupsRepository struct {
	db *sql.DB
}

func NewTaskGroupsRepository(db *sql.DB) *TaskGroupsRepository {
	return &TaskGroupsRepository{db: db}
}
func (r *TaskGroupsRepository) Create(ctx context.Context, taskGroup *models.TaskGroup) error {
	result, err := r.db.ExecContext(ctx,
		"INSERT INTO task_groups (name) VALUES (?)",
		taskGroup.Name)
	if err != nil {
		return err
	}
	taskGroup.ID, _ = result.LastInsertId()
	return nil
}

func (r *TaskGroupsRepository) List(ctx context.Context) ([]models.TaskGroup, error) {
	rows, err := r.db.QueryContext(ctx,
		"SELECT id,name FROM task_groups")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var taskGroups []models.TaskGroup

	for rows.Next() {
		var t models.TaskGroup

		if err := rows.Scan(&t.ID, &t.Name); err != nil {
			return nil, err
		}
		taskGroups = append(taskGroups, t)
	}
	return taskGroups, nil
}

func (r *TaskGroupsRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM task_groups WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func (r *TaskGroupsRepository) GetById(ctx context.Context, id int64) (*models.TaskGroup, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id,name FROM task_groups WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var taskGroup models.TaskGroup

	if err := row.Scan(&taskGroup.ID, &taskGroup.Name); err != nil {
		return nil, err
	}
	return &taskGroup, nil
}

func (r *TaskGroupsRepository) GetByName(ctx context.Context, name string) (*models.TaskGroup, error) {
	row := r.db.QueryRowContext(ctx, "SELECT id,name FROM task_groups WHERE name = ?", name)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var taskGroup models.TaskGroup

	if err := row.Scan(&taskGroup.ID, &taskGroup.Name); err != nil {
		return nil, err
	}
	return &taskGroup, nil
}
