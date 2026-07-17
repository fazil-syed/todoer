package task

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) completeTaskCommand(ctx context.Context, cmd *cli.Command) error {

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
	if task.Status == "DONE" {
		return errors.New("Task already in DONE state")
	}
	if err := c.tasksRepository.Complete(ctx, int64(id)); err != nil {
		return err
	}
	now := time.Now()
	if !task.StartedAt.Valid {
		if err := c.tasksRepository.UpdateStartedAtTime(ctx, int64(id), &now); err != nil {
			return err
		}
	}
	if err := c.tasksRepository.UpdateCompletedAtTime(ctx, int64(id), &now); err != nil {
		return err
	}
	fmt.Println("completed task", task.Title)
	return nil

}

func (c *TaskCommand) CompletTaskCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "complete",
		Aliases: []string{"c", "mark-done", "md"},
		Usage:   "Complete a task",
		Action:  c.completeTaskCommand,
	}
	return cmd
}
