package commands

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *Commands) listTasksHandler(ctx context.Context, cmd *cli.Command) error {
	tasks, err := c.tasksRepository.List(ctx)
	if err != nil {
		return err
	}
	// Print the first line

	fmt.Println("completed	task	id ")

	if len(tasks) == 0 {
		fmt.Println("No tasks found")
	}

	for _, task := range tasks {
		var firstPart string = "[]"
		if task.Done {
			firstPart = "[x]"
		}
		fmt.Println(firstPart, task.Title, task.ID)
	}
	return nil

}

func (c *Commands) ListTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List all tasks",
		Action:  c.listTasksHandler,
	}
	return cmd
}
