package commands

import (
	"context"
	"encoding/csv"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func (c *Commands) exportTasksCommand(ctx context.Context, cmd *cli.Command) error {

	tasks, err := c.tasksRepository.List(ctx)
	if err != nil {
		return err
	}

	file, err := os.Create("tasks.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write header
	if err := writer.Write([]string{"task", "completed", "id"}); err != nil {
		return err
	}
	for _, task := range tasks {
		if err := writer.Write([]string{task.Title, strconv.FormatBool(task.Done), strconv.FormatInt(task.ID, 10)}); err != nil {
			return err
		}
	}
	return nil

}

func (c *Commands) ExportTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "export",
		Aliases: []string{"c"},
		Usage:   "export all tasks to csv",
		Action:  c.exportTasksCommand,
	}
	return cmd
}
