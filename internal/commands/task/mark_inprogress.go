package task

import (
	"context"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) markInprogressTaskCommand(ctx context.Context, cmd *cli.Command) error {

	taskID := cmd.Args().First()

	if taskID == "" {
		fmt.Println("missing task id")
		return nil
	}

	id, err := strconv.Atoi(taskID)
	if err != nil {
		return err
	}
	if err := c.tasksRepository.UpdateStatus(ctx, int64(id), "IN_PROGRESS"); err != nil {
		return err
	}
	task, err := c.tasksRepository.GetById(ctx, int64(id))

	if err != nil {
		return err
	}

	fmt.Println("marked task", task.Title, "as in progress.")
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
