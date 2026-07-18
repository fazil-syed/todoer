package task

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) deleteTasksHandler(ctx context.Context, cmd *cli.Command) error {
	taskID := cmd.Args().First()

	if taskID == "" {
		fmt.Println("missing task id")
		return nil
	}
	id, err := strconv.Atoi(taskID)
	if err != nil {
		return err
	}
	task, err := c.tasksRepository.GetById(ctx, int64(id))

	if err != nil {
		return err
	}
	if err := c.tasksRepository.Delete(ctx, int64(id)); err != nil {
		return err
	}

	fmt.Println("Deleted task")
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintSingleTaskHeadLine()
	printer.PrintSingleTask(*task, false)
	return nil
}

func (c *TaskCommand) DeleteTaskCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "delete",
		Aliases: []string{"d"},
		Usage:   "delete a task",
		Action:  c.deleteTasksHandler,
	}
	return cmd
}
