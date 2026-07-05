package commands

import (
	"github.com/fazil-syed/todoer/internal/repository"
)

type Commands struct {
	tasksRepository *repository.TaskRepository
}

func New(tasksRepository *repository.TaskRepository) *Commands {

	return &Commands{tasksRepository: tasksRepository}

}
