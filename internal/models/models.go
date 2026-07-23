package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int64          `json:"id"`
	Title       string         `json:"title"`
	Status      string         `json:"status"`
	Archived    bool           `json:"archived"`
	Note        sql.NullString `json:"note"`
	GroupId     int64          `json:"group_id"`
	CreatedAt   time.Time      `json:"created_at"`
	StartedAt   sql.NullTime   `json:"started_at"`
	CompletedAt sql.NullTime   `json:"completed_at"`
	GroupName   string         `json:"-"`
}

type TaskGroup struct {
	ID   int64
	Name string
}
