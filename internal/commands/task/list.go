package task

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/urfave/cli/v3"
)

func wrapText(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	line := words[0]
	for _, word := range words[1:] {
		if len(line)+1+(len(word)) <= width {
			line += " " + word
		} else {
			lines = append(lines, line)
			line = word
		}
	}
	lines = append(lines, line)
	return lines
}

func (c *TaskCommand) listAllGroupTasks(ctx context.Context, cmd *cli.Command) error {
	sortOrder := cmd.String("sort")
	tasks, err := c.tasksRepository.GetAllTasksByGroup(ctx, sortOrder)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print the first line
	fmt.Fprintln(w, "Status\tTask\tID\tGroup")
	fmt.Fprintln(w, "----------------------------------------")
	for _, task := range tasks {
		var status string
		switch task.Status {
		case "DONE":
			status = "[x]"
		case "IN_PROGRESS":
			status = "[i]"
		default:
			status = "[ ]"
		}

		lines := wrapText(task.Title, 40)
		// First line
		fmt.Fprintf(w, "%s\t%s\t%d\t%s\n", status, lines[0], task.ID, task.GroupName)

		// Remaining lines
		for _, line := range lines[1:] {
			fmt.Fprintf(w, "\t%s\t\n", line)
		}
	}
	return nil
}

func (c *TaskCommand) listTasksHandler(ctx context.Context, cmd *cli.Command) error {
	groupName := cmd.String("group")
	sortOrder := cmd.String("sort")

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
	tasks, err := c.tasksRepository.List(ctx, taskGroup.ID, sortOrder)
	if err != nil {
		return err
	}
	if len(tasks) == 0 {
		fmt.Println("No tasks found")
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print the first line
	fmt.Fprintln(w, "Status\tTask\tID")
	fmt.Fprintln(w, "----------------------------------------")
	for _, task := range tasks {
		var status string
		switch task.Status {
		case "DONE":
			status = "[x]"
		case "IN_PROGRESS":
			status = "[i]"
		default:
			status = "[ ]"
		}

		lines := wrapText(task.Title, 40)
		// First line
		fmt.Fprintf(w, "%s\t%s\t%d\n", status, lines[0], task.ID)

		// Remaining lines
		for _, line := range lines[1:] {
			fmt.Fprintf(w, "\t%s\t\n", line)
		}
	}
	return nil

}

// func (c *TaskCommand) ListTasksCommand() *cli.Command {
// 	cmd := &cli.Command{
// 		Name:    "list",
// 		Aliases: []string{"l"},
// 		Usage:   "List all tasks",
// 		Action:  c.listTasksHandler,
// 		Flags: []cli.Flag{
// 			&cli.StringFlag{
// 				Name:    "group",
// 				Value:   "default",
// 				Aliases: []string{"g"},
// 				Usage:   "specify which group the task belongs to ",
// 			},
// 			&cli.StringFlag{
// 				Name:    "sort",
// 				Value:   "id",
// 				Aliases: []string{"s"},
// 				Usage:   "sort order for sorting the tasks",
// 			},
// 		},
// 	}
// 	return cmd
// }

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
		},
	}
}
