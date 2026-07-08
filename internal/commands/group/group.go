package group

import (
	"github.com/fazil-syed/todoer/internal/repository"
	"github.com/urfave/cli/v3"
)

type GroupCommand struct {
	taskGroupsRepository *repository.TaskGroupsRepository
}

func NewGroupCommand(taskGroupsRepository *repository.TaskGroupsRepository) *GroupCommand {

	return &GroupCommand{taskGroupsRepository: taskGroupsRepository}

}

func (g *GroupCommand) RegisterGroupCommands() *cli.Command {
	return &cli.Command{
		Name:    "group",
		Aliases: []string{"g"},
		Usage:   "Manage group",
		Commands: []*cli.Command{
			g.AddGroupCommand(),
			g.ListGroupsCommand(),
		},
	}
}
