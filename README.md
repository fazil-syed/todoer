# todoer

`todoer` is a local, SQLite-backed command-line task manager written in Go. It supports the complete day-to-day task workflow: organizing work into groups, tracking TODO, in-progress, and completed states, recording progress timestamps, reviewing tasks across groups, cleaning up completed work, and exporting task data. No account or network connection is required.

## Features

- Create and list task groups, including the automatically created `default` group
- Add tasks, move them between TODO, in-progress, and completed states, and record timestamps
- List a single group's tasks or every task across groups, with sortable output
- Delete individual tasks, clear completed tasks from a group, or purge all tasks with an explicit confirmation flag
- Export group data as CSV or JSON

## Requirements

- Go 1.26.3 or later

## Install

### Homebrew

```sh
brew tap fazil-syed/tap
brew trust fazil-syed/tap
brew install todoer
```

### From source

Clone the repository and build the binary:

```sh
git clone https://github.com/fazil-syed/todoer.git
cd todoer
go build -o todoer ./cmd/todoer
```

Run it from the repository with `./todoer`, or move the binary to a directory on your `PATH` to use `todoer` anywhere.

For development, you can run it without building first:

```sh
go run ./cmd/todoer --help
```

## Quick start

Every new database includes a `default` group.

```sh
# Add a task to the default group
todoer task add Buy milk

# Create a group and add a task to it
todoer group add work
todoer task add Fix login bug --group work

# View tasks and update their state
todoer task list --group work
todoer task mark-inprogress 2
todoer task complete 2
```

Use `todoer --help`, `todoer task --help`, or `todoer group --help` for the CLI's built-in help.

## Commands

### Groups

| Command | Aliases | Description |
| --- | --- | --- |
| `todoer group add <name>` | `group` → `g`; `add` → `a` | Create a task group. |
| `todoer group list` | `group` → `g`; `list` → `l` | List all task groups. |

### Tasks

| Command | Aliases | Description |
| --- | --- | --- |
| `todoer task add <title> [--group <name>]` | `task` → `t`; `add` → `a` | Add a TODO task. The group defaults to `default`. |
| `todoer task list [--group <name>\|all] [--sort id\|created_at\|done]` | `task` → `t`; `list` → `l` | List tasks in one group, or all groups. |
| `todoer task mark-inprogress <id>` | `task` → `t`; `mark-inprogress` → `mi` | Mark a task in progress and set its start time. |
| `todoer task mark-todo <id>` | `task` → `t`; `mark-todo` → `mt` | Mark a task as TODO. |
| `todoer task complete <id>` | `task` → `t`; `complete` → `c`, `mark-done`, `md` | Mark a task done and set its completion time. |
| `todoer task delete <id>` | `task` → `t`; `delete` → `d` | Permanently delete one task. |
| `todoer task clear [--group <name>]` | `task` → `t`; `clear` → `cl` | Permanently remove completed tasks from a group. |
| `todoer task export [--format csv\|json] [--group <name>] [--sort id\|created_at\|done]` | `task` → `t`; `export` → `e` | Export a group's tasks. |
| `todoer task purge --force` | `task` → `t`; `purge` → `p` | Permanently remove all tasks in every group. |

### Flags

| Command | Flag | Alias | Default | Description |
| --- | --- | --- | --- | --- |
| `task add` | `--group <name>` | `-g` | `default` | Assign the new task to a group. |
| `task list` | `--group <name\|all>` | `-g` | `default` | List one group, or use `all` to list every group. |
| `task list` | `--sort id\|created_at\|done` | `-s` | `id` | Set the task sort order. |
| `task clear` | `--group <name>` | `-g` | `default` | Clear completed tasks only from this group. |
| `task export` | `--format csv\|json` | `-f` | `csv` | Choose the export file format. |
| `task export` | `--group <name>` | `-g` | `default` | Export tasks from this group. |
| `task export` | `--sort id\|created_at\|done` | `-s` | `id` | Set the exported task sort order. |
| `task purge` | `--force` | — | off | Required before permanently deleting all tasks. |

## Examples

```sh
# List all tasks, grouped in the output by their group name
todoer task list --group all

# Show a group ordered by creation time
todoer task list -g personal -s created_at

# Clear only completed tasks from the default group
todoer task clear

# Export the work group; creates or overwrites tasks.json in the current directory
todoer task export -g work -f json
```

## Data and exports

`todoer` creates its database automatically and applies bundled migrations on startup. The database is stored at:

| Platform | Location |
| --- | --- |
| macOS | `~/Library/Application Support/todoer/todoer.db` |
| Linux | `~/.local/share/todoer/todoer.db` |
| Windows | `%LOCALAPPDATA%\\todoer\\todoer.db` |

Exports are written to the current working directory as `tasks.csv` or `tasks.json`. Exporting in the same directory replaces the existing file of that format.

## Development

```sh
go build ./cmd/todoer
```

There is currently no automated test suite in this repository.

## License

Licensed under the [Apache License 2.0](LICENSE.txt).
