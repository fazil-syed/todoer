package main

import (
	"context"

	"log"

	"github.com/fazil-syed/todoer/internal/cli"
	"github.com/fazil-syed/todoer/internal/commands"
	"github.com/fazil-syed/todoer/internal/db"
	"github.com/fazil-syed/todoer/internal/repository"
	_ "modernc.org/sqlite"
)

func main() {

	ctx := context.Background()

	db.Init(ctx)

	defer db.DB.Close()

	db.Migrate(ctx)

	todoer, err := cli.NewTodoer()
	if err != nil {
		log.Fatal(err)
	}

	tasksRepo := repository.NewTaskRepository(db.DB)
	cmds := commands.New(tasksRepo)

	todoer.RegisterCommand(cmds.AddTasksCommand())
	todoer.RegisterCommand(cmds.ListTasksCommand())
	todoer.RegisterCommand(cmds.CompletTaskCommand())

	todoer.StartTodoer(ctx)
}
