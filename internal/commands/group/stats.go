package group

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fazil-syed/todoer/internal/types"
	"github.com/urfave/cli/v3"
)

func (c *GroupCommand) listGroupStatsHandler(ctx context.Context, cmd *cli.Command) error {

	groupName := cmd.String("group")

	var allGroupStats []types.GroupStats
	if groupName == "all" {
		allStats, err := c.taskGroupsRepository.GetAllGroupStats(ctx)
		if err != nil {
			return err
		}
		allGroupStats = allStats
	} else {

		taskGroup, err := c.taskGroupsRepository.GetByName(ctx, groupName)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				println("group not found")
				return nil
			}
			return err
		}
		groupStats, err := c.taskGroupsRepository.GetGroupStats(ctx, taskGroup.Name)
		if err != nil {
			return err
		}
		allGroupStats = append(allGroupStats, *groupStats)
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	defer w.Flush()

	fmt.Fprintf(w, "Group\tTODO\tIN_PROGRESS\tDONE\tTOTAL\n")
	fmt.Fprintf(w, "-----\t----\t-----------\t---\t-----\n")
	for _, stats := range allGroupStats {
		fmt.Fprintf(w, "%s\t%d\t%d\t%d\t%d\n", stats.GroupName, stats.Todo, stats.InProgress, stats.Done, stats.Total)
	}

	return nil

}

func (c *GroupCommand) ListGroupStatsCommand() *cli.Command {
	cmd := &cli.Command{
		Name:    "stats",
		Aliases: []string{"st"},
		Usage:   "Get Stats for group(s)",
		Action:  c.listGroupStatsHandler,
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
