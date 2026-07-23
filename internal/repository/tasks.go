package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

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

func (r *TaskRepository) ListByStatusAndGroup(ctx context.Context, groupID int64, status string, fetchArchived bool) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,title,status,started_at,completed_at,created_at,group_id,note,archived
		 FROM tasks
		WHERE group_id = ? AND status = ?
			AND (? = TRUE OR archived = FALSE)
			ORDER BY
				id DESC
		`, groupID, status, fetchArchived)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.StartedAt, &t.CompletedAt, &t.CreatedAt, &t.GroupId, &t.Note, &t.Archived); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
func (r *TaskRepository) List(ctx context.Context, groupID int64, orderBy string, fetchArchived bool) ([]models.Task, error) {
	fmt.Println(fetchArchived)
	rows, err := r.db.QueryContext(ctx,
		`SELECT id,title,status,started_at,completed_at,created_at,group_id,note,archived
		 FROM tasks
		WHERE group_id = ?
			AND (? = TRUE OR archived = FALSE)
			ORDER BY
				CASE WHEN ? = 'done' THEN status END DESC,
				CASE WHEN ? = 'created_at' THEN created_at END ASC,
				id DESC
		`, groupID, fetchArchived, orderBy, orderBy)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.StartedAt, &t.CompletedAt, &t.CreatedAt, &t.GroupId, &t.Note, &t.Archived); err != nil {
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
func (r *TaskRepository) AddTaskNote(ctx context.Context, id int64, note string) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET note = ? WHERE id = ?", note, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *TaskRepository) UpdateArchiveStatus(ctx context.Context, id int64, archived bool) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET archived = ? WHERE id = ?", archived, id)
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

func (r *TaskRepository) UpdateStartedAtTime(ctx context.Context, id int64, at_time *time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET started_at = ? WHERE id = ?", at_time, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *TaskRepository) UpdateCompletedAtTime(ctx context.Context, id int64, at_time *time.Time) error {
	_, err := r.db.ExecContext(ctx,
		"UPDATE tasks SET completed_at = ? WHERE id = ?", at_time, id)
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
	row := r.db.QueryRowContext(ctx, "SELECT id,title,status,created_at,started_at,completed_at,group_id,note,archived FROM tasks WHERE id = ?", id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var task models.Task

	if err := row.Scan(&task.ID, &task.Title, &task.Status, &task.CreatedAt, &task.StartedAt, &task.CompletedAt, &task.GroupId, &task.Note, &task.Archived); err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GetAllTasksByGroup(ctx context.Context, orderBy string, fetchArchived bool) ([]models.Task, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT t.id,t.title,t.status,t.created_at,t.started_at,t.completed_at,t.group_id,t.note,t.archived,g.name
		 FROM tasks t
		 INNER JOIN task_groups g ON t.group_id=g.id
		 WHERE (? = TRUE OR t.archived = FALSE)
			ORDER BY
				CASE WHEN ? = 'done' THEN status END DESC,
				CASE WHEN ? = 'created_at' THEN created_at END ASC,
				g.name ASC,
				t.id DESC
		`, fetchArchived, orderBy, orderBy)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task

	for rows.Next() {
		var t models.Task

		if err := rows.Scan(&t.ID, &t.Title, &t.Status, &t.CreatedAt, &t.StartedAt, &t.CompletedAt, &t.GroupId, &t.Note, &t.Archived, &t.GroupName); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}
