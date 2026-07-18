package task

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/fazil-syed/todoer/internal/models"
)

type TaskPrinter struct {
	w *tabwriter.Writer
}

func NewTaskPrinter(out io.Writer) *TaskPrinter {

	return &TaskPrinter{

		w: tabwriter.NewWriter(out, 0, 0, 2, ' ', 0),
	}

}

func (tp *TaskPrinter) Flush() error {

	return tp.w.Flush()

}

func wrapText(s string, width int) []string {
	words := strings.Fields(s)
	if len(words) == 0 {
		return []string{""}
	}
	var lines []string
	line := words[0]
	for _, word := range words[1:] {
		if len(line)+1+(len(word)) <= width {
			line += " " + word
		} else {
			lines = append(lines, line)
			line = word
		}
	}
	lines = append(lines, line)
	return lines
}

func (tp *TaskPrinter) PrintSingleTask(task models.Task, printGroup bool) {
	var status string
	switch task.Status {
	case "DONE":
		status = "[x]"
	case "IN_PROGRESS":
		status = "[i]"
	default:
		status = "[ ]"
	}

	lines := wrapText(task.Title, 40)
	started := "-"
	if task.StartedAt.Valid {
		started = task.StartedAt.Time.Format("02 Jan 2006 03:04 PM")
	}
	completed := "-"
	if task.CompletedAt.Valid {
		completed = task.CompletedAt.Time.Format("02 Jan 2006 03:04 PM")
	}
	// First line
	if printGroup {
		fmt.Fprintf(tp.w, "%s\t%d\t%s\t%s\t%s\t%s\n", task.GroupName, task.ID, status, lines[0], started, completed)

	} else {
		fmt.Fprintf(tp.w, "%d\t%s\t%s\t%s\t%s\n", task.ID, status, lines[0], started, completed)
	}

	// Remaining lines
	for _, line := range lines[1:] {
		if printGroup {
			fmt.Fprintf(tp.w, "\t\t\t%s\t\t\n", line)
		} else {
			fmt.Fprintf(tp.w, "\t\t%s\t\t\n", line)
		}
	}
}

func (tp *TaskPrinter) PrintSingleTaskHeadLine() {
	fmt.Fprintln(tp.w, "ID\tStatus\tTask\tStarted At\tCompleted At")
	fmt.Fprintln(tp.w, "--\t------\t----\t----------\t------------")
}
func (tp *TaskPrinter) PrintTaskHeadLineWithGroup() {
	fmt.Fprintln(tp.w, "Group\tID\tStatus\tTask\tStarted At\tCompleted At")
	fmt.Fprintln(tp.w, "-----\t--\t------\t----\t----------\t------------")
}
