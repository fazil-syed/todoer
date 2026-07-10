package models

import "time"

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	GroupId   int64     `json:"group_id"`
	CreatedAt time.Time `json:"-"`
}

type TaskGroup struct {
	ID   int64
	Name string
}
