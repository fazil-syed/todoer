package task

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/fazil-syed/todoer/internal/models"
	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) completeTaskCommand(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args().Slice()

	note := cmd.String("note")

	if len(args) == 0 {
		fmt.Println("missing task id")
		return nil
	}

	ids := make([]int64, 0, len(args))

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil {
			return fmt.Errorf("invalid task id %q: %w", arg, err)
		}
		ids = append(ids, int64(id))
	}
	var tasks []models.Task
	for _, id := range ids {

		task, err := c.tasksRepository.GetById(ctx, int64(id))

		if err != nil {
			return err
		}
		if task.Status == "DONE" {
			return fmt.Errorf("Task %d already in DONE state", task.ID)
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
		if note != "" {
			fmt.Println("storting note")
			if err := c.tasksRepository.AddTaskNote(ctx, int64(id), note); err != nil {
				return err
			}
		}
		task, err = c.tasksRepository.GetById(ctx, int64(task.ID))

		if err != nil {
			return err
		}
		tasks = append(tasks, *task)

	}
	fmt.Println("Completed tasks")
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintSingleTaskHeadLine()
	for _, task := range tasks {
		printer.PrintSingleTask(task, false)
	}
	return nil

}

func (c *TaskCommand) CompletTaskCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "complete",
		Aliases: []string{"c", "mark-done", "md"},
		Usage:   "Complete a task",
		Action:  c.completeTaskCommand,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "note",
				Aliases: []string{"n"},
				Usage:   "store a optional note for task completion",
			},
		},
	}
	return cmd
}
