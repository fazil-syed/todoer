package models

import "time"

type Task struct {
	ID        int64
	Title     string
	Done      bool
	CreatedAt time.Time
}
