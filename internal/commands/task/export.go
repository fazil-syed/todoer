package task

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
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
		if err := writer.Write([]string{task.Title, task.Status, strconv.FormatInt(task.ID, 10)}); err != nil {
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

func (c *TaskCommand) exportTasksCommand(ctx context.Context, cmd *cli.Command) error {

	groupName := cmd.String("group")
	sortOrder := cmd.String("sort")

	taskGroup, err := c.taskGroupsRepository.GetByName(ctx, groupName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			println("group not found")
			return nil
		}
		return err
	}
	switch sortOrder {
	case "done", "created_at", "id":
	default:
		fmt.Println("invalid sort order")
		return nil
	}

	tasks, err := c.tasksRepository.List(ctx, taskGroup.ID, sortOrder)
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

func (c *TaskCommand) ExportTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "export",
		Aliases: []string{"e"},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "format",
				Value:   "csv",
				Aliases: []string{"f"},
				Usage:   "output format for export",
			},
			&cli.StringFlag{
				Name:    "group",
				Value:   "default",
				Aliases: []string{"g"},
				Usage:   "specify which group the task belongs to ",
			},
			&cli.StringFlag{
				Name:    "sort",
				Value:   "id",
				Aliases: []string{"s"},
				Usage:   "sort order for sorting the tasks",
			},
		},

		Usage:  "export all tasks",
		Action: c.exportTasksCommand,
	}
	return cmd
}
