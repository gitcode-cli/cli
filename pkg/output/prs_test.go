package output

import (
	"bytes"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
)

func TestPRListPrinterTable(t *testing.T) {
	printer, err := NewPRListPrinter(PRListOptions{Format: FormatTable})
	if err != nil {
		t.Fatalf("NewPRListPrinter() error = %v", err)
	}

	prs := []api.PullRequest{
		{
			Number: 1,
			State:  "open",
			Title:  "Fix login bug",
			User:   &api.User{Login: "alice"},
			Reviewers: []*api.User{
				{Login: "reviewer1"},
			},
		},
		{
			Number: 2,
			State:  "open",
			Draft:  true,
			Title:  "Refactor CLI output",
			User:   &api.User{Login: "bob"},
		},
	}

	var buf bytes.Buffer
	if err := printer.Print(&buf, prs); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	output := buf.String()
	for _, want := range []string{"NUMBER", "STATE", "AUTHOR", "REVIEW", "Fix login bug", "requested", "draft"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output = %q, want substring %q", output, want)
		}
	}
}
