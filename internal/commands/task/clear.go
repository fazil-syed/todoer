package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) clearTaskHandler(ctx context.Context, cmd *cli.Command) error {
	groupName := cmd.String("group")

	taskGroup, err := c.taskGroupsRepository.GetByName(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			println("group not found")
			return nil
		}
		return err
	}
	tasks, err := c.tasksRepository.ListByStatusAndGroup(ctx, taskGroup.ID, "DONE")
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}
	if err := c.tasksRepository.DeleteCompleted(ctx, taskGroup.ID); err != nil {
		return err
	}
	fmt.Println("completed tasks cleared")
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintSingleTaskHeadLine()
	for _, task := range tasks {
		printer.PrintSingleTask(task, false)
	}
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
				Name:    "group",
				Value:   "default",
				Aliases: []string{"g"},
				Usage:   "specify which group the task belongs to ",
			},
		},
	}
	return cmd
}
