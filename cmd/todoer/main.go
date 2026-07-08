package main

import (
	"context"

	"log"

	"github.com/fazil-syed/todoer/internal/cli"
	"github.com/fazil-syed/todoer/internal/commands/group"
	"github.com/fazil-syed/todoer/internal/commands/task"
	"github.com/fazil-syed/todoer/internal/db"
	"github.com/fazil-syed/todoer/internal/repository"
	_ "modernc.org/sqlite"
)

var version = "dev"

func main() {

	ctx := context.Background()

	database, err := db.Init(ctx)

	if err != nil {
		log.Fatal(err)
	}

	defer database.Close()

	if err := db.Migrate(ctx, database); err != nil {
		log.Fatal(err)
	}

	todoer, err := cli.NewTodoer(version)
	if err != nil {
		log.Fatal(err)
	}

	tasksRepo := repository.NewTaskRepository(database)
	taskGroupsRepo := repository.NewTaskGroupsRepository(database)

	taskCommand := task.NewTaskCommand(tasksRepo, taskGroupsRepo)
	groupCommand := group.NewGroupCommand(taskGroupsRepo)

	todoer.RegisterCommand(taskCommand.RegisterTaskCommands())
	todoer.RegisterCommand(groupCommand.RegisterGroupCommands())

	todoer.StartTodoer(ctx)
}
