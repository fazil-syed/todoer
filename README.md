# todoer

`todoer` is a local, SQLite-backed command-line task manager written in Go. It supports the complete day-to-day task workflow: organizing work into groups, tracking TODO, in-progress, and completed states, recording progress timestamps and completion notes, archiving completed work, reviewing tasks across groups, and exporting task data. No account or network connection is required.

## Features

- Create and list task groups, including the automatically created `default` group
- View TODO, in-progress, completed, and total task counts for one group or all groups
- Add tasks, move them between TODO, in-progress, and completed states, and record timestamps and optional completion notes
- Archive completed tasks and restore them later; archived tasks are hidden from lists and exports by default
- List a single group's tasks or every task across groups, with sortable output
- Delete tasks, archive all completed tasks in a group, or purge all tasks with an explicit confirmation flag
- Export group data as CSV or JSON with a custom output filename



## Install

### Homebrew

```sh
brew tap fazil-syed/tap
brew trust fazil-syed/tap
brew install todoer
```

### From source

#### Requirements

- Go 1.26.3 or later

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
todoer task complete --note "Fixed in the latest release" 2
```

Use `todoer --help`, `todoer task --help`, or `todoer group --help` for the CLI's built-in help.

### Task titles and the shell

Wrap task titles that contain shell-special characters, such as parentheses, in double quotes so the shell passes them to `todoer` as part of the title:

```sh
todoer task add "Call Mom (Sunday)"
```

If a task title begins with `-`, put `--` after all command flags and before the title. This stops the CLI from treating the title as another flag:

```sh
todoer task add --group work -- "- Review follow-up items"
```

## Commands

### Groups

| Command | Aliases | Description |
| --- | --- | --- |
| `todoer group add <name>` | `group` → `g`; `add` → `a` | Create a task group. |
| `todoer group list` | `group` → `g`; `list` → `l` | List all task groups. |
| `todoer group stats [--group <name>\|all]` | `group` → `g`; `stats` → `st` | Show TODO, in-progress, completed, and total task counts. Defaults to all groups. |

### Tasks

| Command | Aliases | Description |
| --- | --- | --- |
| `todoer task add [--group <name>] [--] <title>` | `task` → `t`; `add` → `a` | Add a TODO task. The group defaults to `default`; use `--` before a title that begins with `-`. |
| `todoer task list [--group <name>\|all] [--sort id\|created_at\|done] [--all]` | `task` → `t`; `list` → `l` | List tasks in one group, or all groups. Archived tasks require `--all`. |
| `todoer task mark-inprogress <id>...` | `task` → `t`; `mark-inprogress` → `mi` | Mark one or more tasks in progress and set their start times. |
| `todoer task mark-todo <id>...` | `task` → `t`; `mark-todo` → `mt` | Mark one or more tasks as TODO. |
| `todoer task complete [--note <text>] <id>...` | `task` → `t`; `complete` → `c`, `mark-done`, `md` | Mark one or more tasks done, set completion times, and optionally save the same note on each task. |
| `todoer task archive <id>...` | `task` → `t`; `archive` → `ar` | Archive one or more completed tasks. |
| `todoer task unarchive <id>...` | `task` → `t`; `unarchive` → `uar` | Restore one or more archived tasks to normal lists and exports. |
| `todoer task delete <id>...` | `task` → `t`; `delete` → `d` | Permanently delete one or more tasks. |
| `todoer task clear [--group <name>]` | `task` → `t`; `clear` → `cl` | Archive all completed tasks in a group. |
| `todoer task export [--format csv\|json] [--group <name>] [--sort id\|created_at\|done] [--output <filename>] [--all]` | `task` → `t`; `export` → `e` | Export a group's tasks. Archived tasks require `--all`. |
| `todoer task purge --force` | `task` → `t`; `purge` → `p` | Permanently remove all tasks in every group. |

### Flags

| Command | Flag | Alias | Default | Description |
| --- | --- | --- | --- | --- |
| `task add` | `--group <name>` | `-g` | `default` | Assign the new task to a group. |
| `group stats` | `--group <name\|all>` | `-g` | `all` | Show statistics for one group, or for every group. |
| `task complete` | `--note <text>` | `-n` | — | Store a completion note. With multiple IDs, the same note is stored on each task. |
| `task list` | `--group <name\|all>` | `-g` | `default` | List one group, or use `all` to list every group. |
| `task list` | `--sort id\|created_at\|done` | `-s` | `id` | Set the task sort order. |
| `task list` | `--all` | — | off | Include archived tasks. |
| `task clear` | `--group <name>` | `-g` | `default` | Archive completed tasks only from this group. |
| `task export` | `--format csv\|json` | `-f` | `csv` | Choose the export file format. |
| `task export` | `--group <name>` | `-g` | `default` | Export tasks from this group. |
| `task export` | `--sort id\|created_at\|done` | `-s` | `id` | Set the exported task sort order. |
| `task export` | `--output <filename>` | `-o` | `tasks` | Set the output filename without its extension. |
| `task export` | `--all` | — | off | Include archived tasks in the export. |
| `task purge` | `--force` | — | off | Required before permanently deleting all tasks. |

## Examples

```sh
# List all tasks, grouped in the output by their group name
todoer task list --group all

# Show a group ordered by creation time
todoer task list -g personal -s created_at

# Show task status counts for every group, or for one group
todoer group stats
todoer group stats -g work

# Complete several tasks and store the same completion note on each
todoer task complete -n "Released to production" 3 4

# Archive a completed task; only completed tasks can be archived
todoer task archive 3

# Include archived tasks when listing, then restore a task to normal lists
todoer task list --all
todoer task unarchive 3

# Archive every completed task from the default group
todoer task clear

# Export the work group, including archived tasks, to work-tasks.json
todoer task export -g work -f json -o work-tasks --all
```

## Data and exports

`todoer` creates its database automatically and applies bundled migrations on startup. The database is stored at:

| Platform | Location |
| --- | --- |
| macOS | `~/Library/Application Support/todoer/todoer.db` |
| Linux | `~/.local/share/todoer/todoer.db` |
| Windows | `%LOCALAPPDATA%\\todoer\\todoer.db` |

Archived tasks are excluded from lists and exports by default. Use `--all` with `task list` or `task export` to include both active and archived tasks.

Exports are written to the current working directory. By default, the files are named `tasks.csv` or `tasks.json`; use `--output` to supply a different base filename. Exporting with the same filename and format replaces the existing file. JSON exports include task metadata, including archive status and completion notes; CSV exports retain their existing task, status, and timestamp columns.

## Development

```sh
go build ./cmd/todoer
```

There is currently no automated test suite in this repository.

## License

Licensed under the [Apache License 2.0](LICENSE.txt).
