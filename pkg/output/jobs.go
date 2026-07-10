package output

import (
	"fmt"
	"io"

	"gitcode.com/gitcode-cli/cli/api"
	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

// WorkflowJobListOptions configures workflow job list text output.
type WorkflowJobListOptions struct {
	Format Format
	Color  *iostreams.ColorScheme
}

// WorkflowJobListPrinter renders workflow job lists.
type WorkflowJobListPrinter struct {
	opts WorkflowJobListOptions
}

// NewWorkflowJobListPrinter validates and returns a workflow job printer.
func NewWorkflowJobListPrinter(opts WorkflowJobListOptions) (*WorkflowJobListPrinter, error) {
	if opts.Format == "" {
		opts.Format = FormatSimple
	}
	return &WorkflowJobListPrinter{opts: opts}, nil
}

// Print renders workflow jobs according to the configured format.
func (p *WorkflowJobListPrinter) Print(w io.Writer, jobs []api.WorkflowRunJob) error {
	switch p.opts.Format {
	case FormatTable:
		return p.printTable(w, jobs)
	default:
		return p.printSimple(w, jobs)
	}
}

func (p *WorkflowJobListPrinter) printSimple(w io.Writer, jobs []api.WorkflowRunJob) error {
	maxNameWidth := len("NAME")
	maxStatusWidth := len("COMPLETED")
	for _, job := range jobs {
		if width := len(jobName(job)); width > maxNameWidth {
			maxNameWidth = width
		}
		if width := len(job.Status); width > maxStatusWidth {
			maxStatusWidth = width
		}
	}

	for _, job := range jobs {
		fmt.Fprintf(
			w,
			"%-*s %s %s\n",
			maxStatusWidth,
			p.statusLabel(job, maxStatusWidth),
			fmt.Sprintf("%-*s", maxNameWidth, jobName(job)),
			stepsLabel(len(job.Steps)),
		)
	}
	return nil
}

func (p *WorkflowJobListPrinter) printTable(w io.Writer, jobs []api.WorkflowRunJob) error {
	maxNameWidth := len("NAME")
	maxIdentifierWidth := len("IDENTIFIER")
	maxSequenceWidth := len("SEQUENCE")
	maxStatusWidth := len("STATUS")

	for _, job := range jobs {
		if width := len(jobName(job)); width > maxNameWidth {
			maxNameWidth = width
		}
		if width := len(job.Identifier); width > maxIdentifierWidth {
			maxIdentifierWidth = width
		}
		if width := fmt.Sprintf("%d", job.Sequence); len(width) > maxSequenceWidth {
			maxSequenceWidth = len(width)
		}
		if width := len(job.Status); width > maxStatusWidth {
			maxStatusWidth = width
		}
	}

	fmt.Fprintf(
		w,
		"%-*s  %-*s  %-*s  %-*s  %s\n",
		maxStatusWidth, "STATUS",
		maxNameWidth, "NAME",
		maxIdentifierWidth, "IDENTIFIER",
		maxSequenceWidth, "SEQUENCE",
		"STEPS",
	)
	for _, job := range jobs {
		fmt.Fprintf(
			w,
			"%-*s  %-*s  %-*s  %-*d  %d\n",
			maxStatusWidth, p.statusLabel(job, maxStatusWidth),
			maxNameWidth, jobName(job),
			maxIdentifierWidth, job.Identifier,
			maxSequenceWidth, job.Sequence,
			len(job.Steps),
		)
	}
	return nil
}

func (p *WorkflowJobListPrinter) statusLabel(job api.WorkflowRunJob, width int) string {
	label := fmt.Sprintf("%-*s", width, job.Status)
	if p.opts.Color == nil {
		return label
	}
	switch job.Status {
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

func jobName(job api.WorkflowRunJob) string {
	if job.Name != "" {
		return job.Name
	}
	return job.Identifier
}

func stepsLabel(n int) string {
	if n == 1 {
		return "1 step"
	}
	return fmt.Sprintf("%d steps", n)
}
