package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/fazil-syed/todoer/internal/models"
	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) addTaskHandler(ctx context.Context, cmd *cli.Command) error {

	taskTitle := strings.Join(cmd.Args().Slice(), " ")

	groupName := cmd.String("group")

	taskGroup, err := c.taskGroupsRepository.GetByName(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			println("group not found")
			return nil
		}
		return err
	}

	task := &models.Task{
		Title:   taskTitle,
		GroupId: taskGroup.ID,
		Done:    false,
	}
	if err := c.tasksRepository.Create(ctx, task); err != nil {
		return err
	}
	fmt.Printf("added task %s to group %s\n", task.Title, taskGroup.Name)
	return nil

}

func (c *TaskCommand) AddTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add a task",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "group",
				Value:   "default",
				Aliases: []string{"g"},
				Usage:   "specify which group the task belongs to ",
			},
		},
		Action: c.addTaskHandler,
		Description: `Examples:

  todoer task add Buy milk

  todoer task add Fix bug --group work

  todoer task add Pay rent --group personal`,
	}
	return cmd
}
