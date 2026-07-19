package repository

import (
	"context"
	"database/sql"

	"github.com/fazil-syed/todoer/internal/models"
	"github.com/fazil-syed/todoer/internal/types"
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

func (r *TaskGroupsRepository) GetGroupStats(ctx context.Context, name string) (*types.GroupStats, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT g.id,g.name,
		SUM(CASE WHEN t.status = 'TODO' THEN 1 ELSE 0 END) AS todo_count,
		SUM(CASE WHEN t.status = 'IN_PROGRESS' THEN 1 ELSE 0 END) AS inprogress_count,
		SUM(CASE WHEN t.status = 'DONE' THEN 1 ELSE 0 END) AS done_count,
		COUNT(t.id) AS total_count
		FROM tasks t
		INNER JOIN task_groups g ON t.group_id=g.id
			WHERE g.name = ?
			GROUP BY g.name,g.id
	`, name)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var groupStats types.GroupStats

	if err := row.Scan(&groupStats.GroupID, &groupStats.GroupName, &groupStats.Todo, &groupStats.InProgress, &groupStats.Done, &groupStats.Total); err != nil {
		return nil, err
	}

	return &groupStats, nil
}

func (r *TaskGroupsRepository) GetAllGroupStats(ctx context.Context) ([]types.GroupStats, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT g.id,g.name,
		SUM(CASE WHEN t.status = 'TODO' THEN 1 ELSE 0 END) AS todo_count,
		SUM(CASE WHEN t.status = 'IN_PROGRESS' THEN 1 ELSE 0 END) AS inprogress_count,
		SUM(CASE WHEN t.status = 'DONE' THEN 1 ELSE 0 END) AS done_count,
		COUNT(t.id) AS total_count
		FROM tasks t
		INNER JOIN task_groups g ON t.group_id=g.id
			GROUP BY g.name,g.id
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var groupStats []types.GroupStats
	for rows.Next() {
		var gs types.GroupStats
		if err := rows.Scan(&gs.GroupID, &gs.GroupName, &gs.Todo, &gs.InProgress, &gs.Done, &gs.Total); err != nil {
			return nil, err
		}
		groupStats = append(groupStats, gs)

	}

	return groupStats, nil
}
