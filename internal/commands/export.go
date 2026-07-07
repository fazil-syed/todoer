package commands

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/fazil-syed/todoer/internal/models"
	"github.com/urfave/cli/v3"
)

func exportCsv(tasks []models.Task) error {
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

func exportJson(tasks []models.Task) error {
	file, err := os.Create("tasks.json")
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := json.NewEncoder(file)

	encoder.SetIndent("", " ")

	if err := encoder.Encode(tasks); err != nil {
		return err
	}

	return nil
}

func (c *Commands) exportTasksCommand(ctx context.Context, cmd *cli.Command) error {

	tasks, err := c.tasksRepository.List(ctx)
	if err != nil {
		return err
	}
	format := cmd.String("format")
	switch format {
	case "csv":
		exportCsv(tasks)

	case "json":
		exportJson(tasks)

	default:
		fmt.Println("unknown format specified")

	}

	return nil

}

func (c *Commands) ExportTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "export",
		Aliases: []string{"e"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "format",
				Value: "csv",
				Usage: "output format for export",
			},
		},
		Usage:  "export all tasks to csv",
		Action: c.exportTasksCommand,
	}
	return cmd
}
