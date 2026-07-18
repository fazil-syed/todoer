package task

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) markTodoTaskCommand(ctx context.Context, cmd *cli.Command) error {

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
	if task.Status == "TODO" {
		return errors.New("Task already in TODO state")
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
	fmt.Println("Updated task to Todo")
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintSingleTaskHeadLine()
	printer.PrintSingleTask(*task, false)
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
