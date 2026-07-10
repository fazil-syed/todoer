package task

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *TaskCommand) purgeTaskHandler(ctx context.Context, cmd *cli.Command) error {

	if !cmd.IsSet("force") {
		fmt.Println("please use --force option to purge all tasks, this action cannot be reverted. use only if you are sure you want to delete all tasks")
		return nil
	}

	if err := c.tasksRepository.Truncate(ctx); err != nil {
		return err
	}
	fmt.Println("purged all tasks")
	return nil

}

func (c *TaskCommand) PurgeTasksCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "purge",
		Aliases: []string{"p"},
		Usage:   "Purge all tasks",
		Action:  c.purgeTaskHandler,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "force",
				Usage: "this option is required to purge all tasks from all groups. proceed cautiously",
			},
		},
	}
	return cmd
}
