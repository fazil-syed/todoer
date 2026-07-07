package models

import "time"

type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"completed"`
	CreatedAt time.Time `json:"-"`
}
