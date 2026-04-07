package output

import (
	"bytes"
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/api"
)

func TestRepoListPrinterTable(t *testing.T) {
	printer, err := NewRepoListPrinter(RepoListOptions{Format: FormatTable})
	if err != nil {
		t.Fatalf("NewRepoListPrinter() error = %v", err)
	}

	repos := []api.Repository{
		{
			FullName:    "owner/app",
			Private:     false,
			Language:    "Go",
			Description: "CLI application",
		},
		{
			FullName:    "owner/web",
			Private:     true,
			Language:    "",
			Description: "Web application",
		},
	}

	var buf bytes.Buffer
	if err := printer.Print(&buf, repos); err != nil {
		t.Fatalf("Print() error = %v", err)
	}

	output := buf.String()
	for _, want := range []string{"NAME", "VISIBILITY", "LANGUAGE", "owner/app", "public", "private", "CLI application"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output = %q, want substring %q", output, want)
		}
	}
}
