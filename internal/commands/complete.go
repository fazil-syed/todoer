package commands

import (
	"context"
	"fmt"
	"strconv"

	"github.com/urfave/cli/v3"
)

func (c *Commands) completeTaskCommand(ctx context.Context, cmd *cli.Command) error {

	taskID := cmd.Args().First()
	id, err := strconv.Atoi(taskID)
	if err != nil {
		return err
	}
	if err := c.tasksRepository.Complete(ctx, int64(id)); err != nil {
		return err
	}
	task, err := c.tasksRepository.GetById(ctx, int64(id))

	if err != nil {
		return err
	}

	fmt.Println("completed task", task.Title)
	return nil

}

func (c *Commands) CompletTaskCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "complete",
		Aliases: []string{"c"},
		Usage:   "Complete a task",
		Action:  c.completeTaskCommand,
	}
	return cmd
}
