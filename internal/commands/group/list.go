package group

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

func (c *GroupCommand) listGroupsHandler(ctx context.Context, cmd *cli.Command) error {

	groups, err := c.taskGroupsRepository.List(ctx)
	if err != nil {
		return err
	}

	if len(groups) == 0 {
		fmt.Println("No groups found")
	}
	fmt.Println("ID \t Name")
	for _, group := range groups {
		fmt.Println(group.ID, "\t", group.Name)
	}
	return nil

}

func (c *GroupCommand) ListGroupsCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "list",
		Aliases: []string{"l"},
		Usage:   "List all groups",
		Action:  c.listGroupsHandler,
	}
	return cmd
}
