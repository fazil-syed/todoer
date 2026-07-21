package task

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fazil-syed/todoer/internal/models"
	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) markInprogressTaskCommand(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()

	if len(args) == 0 {
		fmt.Println("missing task id")
		return nil
	}

	ids := make([]int64, 0, len(args))

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("invalid task id %q: %w", arg, err)
		}
		ids = append(ids, int64(id))
	}
	var tasks []models.Task
	for _, id := range ids {

		task, err := c.tasksRepository.GetById(ctx, int64(id))

		if err != nil {
			return err
		}
		if task.Status == "IN_PROGRESS" {
			return fmt.Errorf("Task %d already in IN_PROGRESS state", task.ID)
		}
		if err := c.tasksRepository.UpdateStatus(ctx, int64(id), "IN_PROGRESS"); err != nil {
			return err
		}
		now := time.Now()
		if err := c.tasksRepository.UpdateStartedAtTime(ctx, int64(id), &now); err != nil {
			return err
		}
		if err := c.tasksRepository.UpdateCompletedAtTime(ctx, int64(id), nil); err != nil {
			return err
		}
		task, err = c.tasksRepository.GetById(ctx, int64(task.ID))

		if err != nil {
			return err
		}
		tasks = append(tasks, *task)

	}

	fmt.Println("Updated tasks to inprogress")
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintSingleTaskHeadLine()
	for _, task := range tasks {
		printer.PrintSingleTask(task, false)
	}
	return nil

}

func (c *TaskCommand) MarkInprogressTaskCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "mark-inprogress",
		Aliases: []string{"mi"},
		Usage:   "mark a task as in progress",
		Action:  c.markInprogressTaskCommand,
	}
	return cmd
}
