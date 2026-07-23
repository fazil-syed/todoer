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

func exportCsv(tasks []models.Task, filename string) error {
	file, err := os.Create(fmt.Sprintf("%s.csv", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write header
	if err := writer.Write([]string{"id", "status", "task", "started_at", "completed_at"}); err != nil {
		return err
	}
	for _, task := range tasks {
		started := "-"
		if task.StartedAt.Valid {
			started = task.StartedAt.Time.Format("02 Jan 2006 03:04 PM")
		}
		completed := "-"
		if task.CompletedAt.Valid {
			completed = task.CompletedAt.Time.Format("02 Jan 2006 03:04 PM")
		}
		if err := writer.Write([]string{strconv.FormatInt(task.ID, 10), task.Status, task.Title, started, completed}); err != nil {
			return err
		}
	}
	return nil
}

func exportJson(tasks []models.Task, filename string) error {
	file, err := os.Create(fmt.Sprintf("%s.json", filename))
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
	fileName := cmd.String("output")
	fetchArchived := cmd.Bool("all")
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

	tasks, err := c.tasksRepository.List(ctx, taskGroup.ID, sortOrder, fetchArchived)
	if err != nil {
		return err
	}
	format := cmd.String("format")
	switch format {
	case "csv":
		exportCsv(tasks, fileName)

	case "json":
		exportJson(tasks, fileName)

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
			&cli.StringFlag{
				Name:    "output",
				Value:   "tasks",
				Aliases: []string{"o"},
				Usage:   "output file name for export",
			},
			&cli.BoolFlag{
				Name:  "all",
				Value: false,
				Usage: "list archived tasks as well",
			},
		},

		Usage:  "export all tasks",
		Action: c.exportTasksCommand,
	}
	return cmd
}
