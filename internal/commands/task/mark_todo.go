package task

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/fazil-syed/todoer/internal/models"
	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) markTodoTaskCommand(ctx context.Context, cmd *cli.Command) error {

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
		if task.Status == "TODO" {
			return fmt.Errorf("Task %d already in TODO state", task.ID)
		}
		if err := c.tasksRepository.UpdateStatus(ctx, int64(id), "TODO"); err != nil {
			return err
		}
		if err := c.tasksRepository.UpdateStartedAtTime(ctx, int64(id), nil); err != nil {
			return err
		}
		if err := c.tasksRepository.UpdateCompletedAtTime(ctx, int64(id), nil); err != nil {
			return err
		}

		task, err = c.tasksRepository.GetById(ctx, int64(id))

		if err != nil {
			return err
		}
		tasks = append(tasks, *task)

	}

	fmt.Println("Updated tasks to Todo")
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintSingleTaskHeadLine()
	for _, task := range tasks {
		printer.PrintSingleTask(task, false)
	}
	return nil

}

func (c *TaskCommand) MarkTodoTaskCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "mark-todo",
		Aliases: []string{"mt"},
		Usage:   "mark a task as todo",
		Action:  c.markTodoTaskCommand,
	}
	return cmd
}
