package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *Commands) clearTaskHandler(ctx context.Context, cmd *cli.Command) error {

	if err := c.tasksRepository.DeleteCompleted(ctx); err != nil {
		return err
	}
	fmt.Println("completed tasks cleared")
	return nil

}

func (c *Commands) ClearTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "clear",
		Aliases: []string{"cl"},
		Usage:   "Clear completed tasks",
		Action:  c.clearTaskHandler,
	}
	return cmd
}
