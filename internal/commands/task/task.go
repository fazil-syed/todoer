package task

import (
	"github.com/fazil-syed/todoer/internal/repository"
	"github.com/urfave/cli/v3"
)

type TaskCommand struct {
	tasksRepository      *repository.TaskRepository
	taskGroupsRepository *repository.TaskGroupsRepository
}

func NewTaskCommand(tasksRepository *repository.TaskRepository, taskGroupsRepository *repository.TaskGroupsRepository) *TaskCommand {

	return &TaskCommand{tasksRepository: tasksRepository, taskGroupsRepository: taskGroupsRepository}

}

func (t *TaskCommand) RegisterTaskCommands() *cli.Command {
	return &cli.Command{
		Name:    "task",
		Aliases: []string{"t"},
		Usage:   "Manage tasks",
		Commands: []*cli.Command{
			t.AddTasksCommand(),
			t.ListTasksCommand(),
			t.ClearTasksCommand(),
			t.CompletTaskCommand(),
			t.DeleteTaskCommand(),
			t.ExportTasksCommand(),
			t.PurgeTasksCommand(),
			t.MarkInprogressTaskCommand(),
			t.MarkTodoTaskCommand(),
			t.ArchiveTaskCommand(),
			t.UnArchiveTaskCommand(),
		},
	}
}
