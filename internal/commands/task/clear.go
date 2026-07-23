package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/fazil-syed/todoer/internal/models"
	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) clearTaskHandler(ctx context.Context, cmd *cli.Command) error {
	groupName := cmd.String("group")

	taskGroup, err := c.taskGroupsRepository.GetByName(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			println("group not found")
			return nil
		}
		return err
	}
	tasks, err := c.tasksRepository.ListByStatusAndGroup(ctx, taskGroup.ID, "DONE", true)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}
	var archivedTask []models.Task
	for _, task := range tasks {
		if err := c.tasksRepository.UpdateArchiveStatus(ctx, int64(task.ID), true); err != nil {
			return err
		}
		task, err := c.tasksRepository.GetById(ctx, int64(task.ID))

		if err != nil {
			return err
		}
		archivedTask = append(archivedTask, *task)

	}
	fmt.Println("completed tasks archived")
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintTaskHeadLineWithGroup()
	for _, task := range archivedTask {
		printer.PrintSingleTask(task, true)
	}
	return nil

}

func (c *TaskCommand) ClearTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "clear",
		Aliases: []string{"cl"},
		Usage:   "Clear completed tasks and mark them as archived",
		Action:  c.clearTaskHandler,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "group",
				Value:   "default",
				Aliases: []string{"g"},
				Usage:   "specify which group the task belongs to ",
			},
		},
	}
	return cmd
}
