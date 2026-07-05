package cli

import (
	"context"
	"log"
	"os"

	"github.com/urfave/cli/v3"
)

type TodoerCli struct {
	rootCmd *cli.Command
}

func NewTodoer() (*TodoerCli, error) {
	rootCmd := &cli.Command{
		Name: "todoer",
		Action: func(ctx context.Context, c *cli.Command) error {
			return nil
		},
	}
	t := &TodoerCli{rootCmd: rootCmd}
	return t, nil
}

func (t *TodoerCli) StartTodoer(ctx context.Context) {

	if err := t.rootCmd.Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}

func (t *TodoerCli) RegisterCommand(cmd *cli.Command) {
	t.rootCmd.Commands = append(t.rootCmd.Commands, cmd)
}
