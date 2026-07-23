package task

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/fazil-syed/todoer/internal/models"
	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) unarchiveTasksHandler(ctx context.Context, cmd *cli.Command) error {
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
		if !task.Archived {
			return fmt.Errorf("Task %d not archived to unarchive", task.ID)
		}

		if err != nil {
			return err
		}
		if err := c.tasksRepository.UpdateArchiveStatus(ctx, int64(id), false); err != nil {
			return err
		}
		task, err = c.tasksRepository.GetById(ctx, int64(task.ID))

		if err != nil {
			return err
		}
		tasks = append(tasks, *task)

	}

	fmt.Println("Unarchived tasks")
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintSingleTaskHeadLine()
	for _, task := range tasks {
		printer.PrintSingleTask(task, false)

	}
	return nil
}

func (c *TaskCommand) UnArchiveTaskCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "unarchive",
		Aliases: []string{"uar"},
		Usage:   "unarchive tasks",
		Action:  c.unarchiveTasksHandler,
	}
	return cmd
}
