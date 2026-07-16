package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Status      string       `json:"status"`
	GroupId     int64        `json:"group_id"`
	CreatedAt   time.Time    `json:"-"`
	StartedAt   sql.NullTime `json:"started_at"`
	CompletedAt sql.NullTime `json:"completed_at"`
	GroupName   string       `json:"-"`
}

type TaskGroup struct {
	ID   int64
	Name string
}
