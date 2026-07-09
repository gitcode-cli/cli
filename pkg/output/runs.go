package output

import (
	"fmt"
	"io"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// WorkflowRunListOptions configures workflow run list text output.
type WorkflowRunListOptions struct {
	Format Format
	Color  *iostreams.ColorScheme
}

// WorkflowRunListPrinter renders workflow run lists.
type WorkflowRunListPrinter struct {
	opts WorkflowRunListOptions
}

// NewWorkflowRunListPrinter validates and returns a workflow run printer.
func NewWorkflowRunListPrinter(opts WorkflowRunListOptions) (*WorkflowRunListPrinter, error) {
	if opts.Format == "" {
		opts.Format = FormatSimple
	}
	return &WorkflowRunListPrinter{opts: opts}, nil
}

// Print renders workflow runs according to the configured format.
func (p *WorkflowRunListPrinter) Print(w io.Writer, runs []api.WorkflowRun) error {
	switch p.opts.Format {
	case FormatTable:
		return p.printTable(w, runs)
	default:
		return p.printSimple(w, runs)
	}
}

func (p *WorkflowRunListPrinter) printSimple(w io.Writer, runs []api.WorkflowRun) error {
	maxNumWidth := len("#")
	statusWidth := len("COMPLETED")
	for _, run := range runs {
		if width := len(fmt.Sprintf("#%d", run.RunNumber)); width > maxNumWidth {
			maxNumWidth = width
		}
		if width := len(run.Status); width > statusWidth {
			statusWidth = width
		}
	}

	for _, run := range runs {
		fmt.Fprintf(
			w,
			"%-*s %s %s\n",
			maxNumWidth,
			fmt.Sprintf("#%d", run.RunNumber),
			p.statusLabel(run, statusWidth),
			runTitle(run),
		)
	}
	return nil
}

func (p *WorkflowRunListPrinter) printTable(w io.Writer, runs []api.WorkflowRun) error {
	maxNumWidth := len("NUMBER")
	maxStatusWidth := len("STATUS")
	maxEventWidth := len("EVENT")
	maxBranchWidth := len("BRANCH")
	maxWorkflowWidth := len("WORKFLOW")
	maxActorWidth := len("ACTOR")

	for _, run := range runs {
		if width := len(fmt.Sprintf("#%d", run.RunNumber)); width > maxNumWidth {
			maxNumWidth = width
		}
		if width := len(run.Status); width > maxStatusWidth {
			maxStatusWidth = width
		}
		if width := len(run.Event); width > maxEventWidth {
			maxEventWidth = width
		}
		if width := len(run.HeadBranch); width > maxBranchWidth {
			maxBranchWidth = width
		}
		if width := len(run.WorkflowName); width > maxWorkflowWidth {
			maxWorkflowWidth = width
		}
		if width := len(runActor(run)); width > maxActorWidth {
			maxActorWidth = width
		}
	}

	fmt.Fprintf(
		w,
		"%-*s  %-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
		maxNumWidth, "NUMBER",
		maxStatusWidth, "STATUS",
		maxEventWidth, "EVENT",
		maxBranchWidth, "BRANCH",
		maxWorkflowWidth, "WORKFLOW",
		maxActorWidth, "ACTOR",
		"TITLE",
	)
	for _, run := range runs {
		fmt.Fprintf(
			w,
			"%-*s  %-*s  %-*s  %-*s  %-*s  %-*s  %s\n",
			maxNumWidth, fmt.Sprintf("#%d", run.RunNumber),
			maxStatusWidth, p.statusLabel(run, maxStatusWidth),
			maxEventWidth, run.Event,
			maxBranchWidth, run.HeadBranch,
			maxWorkflowWidth, run.WorkflowName,
			maxActorWidth, runActor(run),
			runTitle(run),
		)
	}

	return nil
}

func (p *WorkflowRunListPrinter) statusLabel(run api.WorkflowRun, width int) string {
	label := fmt.Sprintf("%-*s", width, run.Status)
	if p.opts.Color == nil {
		return label
	}
	switch run.Status {
	case "COMPLETED":
		return p.opts.Color.Green(label)
	case "FAILED":
		return p.opts.Color.Red(label)
	case "RUNNING":
		return p.opts.Color.Yellow(label)
	case "CANCELED", "IGNORED", "PAUSED", "SUSPEND":
		return p.opts.Color.Gray(label)
	default:
		return label
	}
}

func runActor(run api.WorkflowRun) string {
	if run.Actor == nil || run.Actor.Login == "" {
		return "-"
	}
	return run.Actor.Login
}

func runTitle(run api.WorkflowRun) string {
	if run.Title != "" {
		return run.Title
	}
	return run.WorkflowName
}
