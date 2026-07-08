package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/fazil-syed/todoer/internal/models"
	"github.com/urfave/cli/v3"
)

func (c *GroupCommand) addGoupHandler(ctx context.Context, cmd *cli.Command) error {
	groupName := strings.Join(cmd.Args().Slice(), " ")

	taskGroup, err := c.taskGroupsRepository.GetByName(ctx, groupName)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if taskGroup != nil {
		fmt.Printf("%s group already exists!\n", groupName)
		return nil
	}
	taskGroup = &models.TaskGroup{
		Name: groupName,
	}
	if err := c.taskGroupsRepository.Create(ctx, taskGroup); err != nil {
		return err
	}
	fmt.Printf("%s group created\n", taskGroup.Name)
	return nil
}

func (c *GroupCommand) AddGroupCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "add",
		Aliases: []string{"a"},
		Usage:   "Add a group",

		Action: c.addGoupHandler,
		Description: `Examples:

  todoer group add work`,
	}
	return cmd
}
