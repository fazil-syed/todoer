package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) listTasksHandler(ctx context.Context, cmd *cli.Command) error {
	groupName := cmd.String("group")

	taskGroup, err := c.taskGroupsRepository.GetByName(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			println("group not found")
		}
		return err
	}

	tasks, err := c.tasksRepository.List(ctx, taskGroup.ID)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	// Print the first line
	fmt.Println("completed	task	id ")

	for _, task := range tasks {
		var firstPart string = "[]"
		if task.Done {
			firstPart = "[x]"
		}
		fmt.Println(firstPart, task.Title, task.ID)
	}
	return nil

}

func (c *TaskCommand) ListTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List all tasks",
		Action:  c.listTasksHandler,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "group",
				Value:   "default",
				Aliases: []string{"g"},
				Usage:   "specify which group the task belongs to ",
			},
		},
	}
	return cmd
}
