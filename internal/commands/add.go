package commands

import (
	"context"
	"fmt"
	"strings"

	"github.com/fazil-syed/todoer/models"
	"github.com/urfave/cli/v3"
)

func (c *Commands) addTaskHandler(ctx context.Context, cmd *cli.Command) error {

	taskTitle := strings.Join(cmd.Args().Slice(), " ")
	task := &models.Task{
		Title: taskTitle,
		Done:  false,
	}
	if err := c.tasksRepository.Create(ctx, task); err != nil {
		return err
	}
	fmt.Println("added task", task.Title)
	return nil

}

func (c *Commands) AddCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add a task",
		Action:  c.addTaskHandler,
	}
	return cmd
}
