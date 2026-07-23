package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) listAllGroupTasks(ctx context.Context, cmd *cli.Command) error {
	sortOrder := cmd.String("sort")
	fetchArchived := cmd.Bool("all")
	tasks, err := c.tasksRepository.GetAllTasksByGroup(ctx, sortOrder, fetchArchived)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}
	// Print the first line
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintTaskHeadLineWithGroup()
	for _, task := range tasks {
		printer.PrintSingleTask(task, true)
	}
	return nil
}

func (c *TaskCommand) listTasksHandler(ctx context.Context, cmd *cli.Command) error {
	groupName := cmd.String("group")
	sortOrder := cmd.String("sort")
	fetchArchived := cmd.Bool("all")

	switch sortOrder {
	case "done", "created_at", "id":
	default:
		fmt.Println("invalid sort order")
		return nil
	}
	if groupName == "all" {
		return c.listAllGroupTasks(ctx, cmd)
	}
	taskGroup, err := c.taskGroupsRepository.GetByName(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			fmt.Println("group not found")
		}
		return err
	}

	tasks, err := c.tasksRepository.List(ctx, taskGroup.ID, sortOrder, fetchArchived)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}
	// Print the first line
	printer := NewTaskPrinter(os.Stdout)
	defer printer.Flush()
	printer.PrintSingleTaskHeadLine()
	for _, task := range tasks {
		printer.PrintSingleTask(task, false)
	}
	return nil

}

func (c *TaskCommand) ListTasksCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List tasks",
		UsageText: `todoer task list [options]

Examples:
  todoer task list
  todoer task list --group work
  todoer task list --group all
  todoer task list --sort created_at
  todoer task list -g work -s done`,
		Action: c.listTasksHandler,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "group",
				Aliases: []string{"g"},
				Value:   "default",
				Usage:   "task group to list (default|all|<group-name>)",
			},
			&cli.StringFlag{
				Name:    "sort",
				Aliases: []string{"s"},
				Value:   "id",
				Usage:   "sort by: id, created_at, or done",
			},
			&cli.BoolFlag{
				Name:  "all",
				Value: false,
				Usage: "list archived tasks as well",
			},
		},
	}
}
