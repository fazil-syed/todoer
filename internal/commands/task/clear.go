package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) clearTaskHandler(ctx context.Context, cmd *cli.Command) error {
	groupName := cmd.String("group")

	taskGroup, err := c.taskGroupsRepository.GetByName(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			println("group not found")
		}
		return err
	}
	if err := c.tasksRepository.DeleteCompleted(ctx, taskGroup.ID); err != nil {
		return err
	}
	fmt.Println("completed tasks cleared")
	return nil

}

func (c *TaskCommand) ClearTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "clear",
		Aliases: []string{"cl"},
		Usage:   "Clear completed tasks",
		Action:  c.clearTaskHandler,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "group",
				Value: "default",
				Usage: "specify which group the task belongs to ",
			},
		},
	}
	return cmd
}
