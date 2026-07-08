package task

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) purgeTaskHandler(ctx context.Context, cmd *cli.Command) error {

	if err := c.tasksRepository.Truncate(ctx); err != nil {
		return err
	}
	fmt.Println("purged all tasks")
	return nil

}

func (c *TaskCommand) PurgeTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "purge",
		Aliases: []string{"p"},
		Usage:   "Purge all tasks",
		Action:  c.purgeTaskHandler,
	}
	return cmd
}
