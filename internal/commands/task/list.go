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
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	defer w.Flush()

	// Print the first line
	fmt.Fprintln(w, "Status\tTask\tID")
	for _, task := range tasks {
		status := "[ ]"
		if task.Done {
			status = "[x]"
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
