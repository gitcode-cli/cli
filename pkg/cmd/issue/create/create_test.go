package create

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdCreate(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "create with title",
			args:    []string{"--title", "Test Issue"},
			wantErr: false,
		},
		{
			name:    "create with title and body",
			args:    []string{"--title", "Test", "--body", "Description"},
			wantErr: false,
		},
		{
			name:    "create with labels",
			args:    []string{"--title", "Test", "--label", "bug,enhancement"},
			wantErr: false,
		},
		{
			name:    "create with template path",
			args:    []string{"--title", "Test", "--template-path", ".gitcode/ISSUE_TEMPLATE/feature.yaml"},
			wantErr: false,
		},
		{
			name:    "create with custom fields json",
			args:    []string{"--title", "Test", "--custom-fields-json", `[{"id":"field","value":"demo"}]`},
			wantErr: false,
		},
		{
			name:    "no title",
			args:    []string{},
			wantErr: false, // Command runs, error in run function
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdCreate(f, func(opts *CreateOptions) error {
				return nil
			})
			cmd.SetArgs(tt.args)

			err := cmd.Execute()
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateRunFailsWhenAssigneesAreNotApplied(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &CreateOptions{
		IO: f.IOStreams,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					switch req.URL.Path {
					case "/api/v5/users/alice":
						return issueResponse(http.StatusOK, `{"id":"101","login":"alice"}`), nil
					case "/api/v5/repos/owner/repo/issues":
						return issueResponse(http.StatusOK, `{"number":"12","html_url":"https://gitcode.com/owner/repo/issues/12"}`), nil
					case "/api/v5/repos/owner/repo/issues/12":
						return issueResponse(http.StatusOK, `{"number":"12","assignees":[]}`), nil
					default:
						t.Fatalf("unexpected request: %s", req.URL.Path)
						return nil, nil
					}
				}),
			}, nil
		},
		Repository: "owner/repo",
		Title:      "Bug report",
		Assignees:  []string{"alice"},
	}

	err := createRun(opts)
	if err == nil {
		t.Fatal("createRun() error = nil, want assignee verification error")
	}
	if !strings.Contains(f.IOStreams.Out.(*bytes.Buffer).String(), "Created issue #12") {
		t.Fatalf("stdout = %q, want created issue output", f.IOStreams.Out.(*bytes.Buffer).String())
	}
	if !strings.Contains(err.Error(), "https://gitcode.com/owner/repo/issues/12") {
		t.Fatalf("createRun() error = %v", err)
	}
	if !strings.Contains(err.Error(), "did not apply the requested assignees") {
		t.Fatalf("createRun() error = %v", err)
	}
}

func TestCreateRunUsesOwnerIssueCreateWhenAdvancedFieldsAreSet(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")

	f := cmdutil.TestFactory()
	opts := &CreateOptions{
		IO:               f.IOStreams,
		Repository:       "owner/repo",
		Title:            "Feature request",
		TemplatePath:     ".gitcode/ISSUE_TEMPLATE/feature.yaml",
		SecurityHole:     true,
		IssueType:        "需求",
		IssueSeverity:    "高",
		CustomFieldsJSON: `[{"id":"field","value":"demo"}]`,
		HttpClient: func() (*http.Client, error) {
			return &http.Client{
				Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
					if req.URL.Path != "/api/v5/repos/owner/issues" {
						t.Fatalf("request path = %s, want /api/v5/repos/owner/issues", req.URL.Path)
					}
					if got := req.Header.Get("Content-Type"); got != "application/json" {
						t.Fatalf("content-type = %q, want application/json", got)
					}
					body, err := io.ReadAll(req.Body)
					if err != nil {
						t.Fatalf("ReadAll() error = %v", err)
					}
					text := string(body)
					for _, want := range []string{
						`"repo":"repo"`,
						`"title":"Feature request"`,
						`"template_path":".gitcode/ISSUE_TEMPLATE/feature.yaml"`,
						`"security_hole":"true"`,
						`"issue_type":"需求"`,
						`"issue_severity":"高"`,
						`"custom_fields":[{"id":"field","value":"demo"}]`,
					} {
						if !strings.Contains(text, want) {
							t.Fatalf("request body = %s, want substring %s", text, want)
						}
					}
					return issueResponse(http.StatusOK, `{"number":"34","html_url":"https://gitcode.com/owner/repo/issues/34"}`), nil
				}),
			}, nil
		},
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}
}

func TestCreateRunDryRunShowsAdvancedFields(t *testing.T) {
	f := cmdutil.TestFactory()
	opts := &CreateOptions{
		IO:               f.IOStreams,
		BaseRepo:         func() (string, error) { return "owner/repo", nil },
		Title:            "Feature request",
		DryRun:           true,
		TemplatePath:     ".gitcode/ISSUE_TEMPLATE/feature.yaml",
		SecurityHole:     true,
		IssueType:        "需求",
		IssueSeverity:    "高",
		CustomFieldsJSON: `[{"id":"field","value":"demo"}]`,
	}

	if err := createRun(opts); err != nil {
		t.Fatalf("createRun() error = %v", err)
	}

	out := f.IOStreams.Out.(*bytes.Buffer).String()
	for _, want := range []string{
		`Dry run: would create issue "Feature request" in owner/repo`,
		`template-path: .gitcode/ISSUE_TEMPLATE/feature.yaml`,
		`security-hole: true`,
		`issue-type: 需求`,
		`issue-severity: 高`,
		`custom-fields: 1 item(s)`,
	} {
		if !strings.Contains(out, want) {
			t.Fatalf("stdout = %q, want substring %q", out, want)
		}
	}
}

func TestGetCustomFields(t *testing.T) {
	tempDir := t.TempDir()
	fieldsPath := filepath.Join(tempDir, "fields.json")
	if err := os.WriteFile(fieldsPath, []byte(`[{"id":"field","value":"demo"}]`), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	tests := []struct {
		name    string
		opts    *CreateOptions
		wantLen int
		wantErr string
	}{
		{
			name:    "from json flag",
			opts:    &CreateOptions{CustomFieldsJSON: `[{"id":"field","value":"demo"}]`},
			wantLen: 1,
		},
		{
			name:    "from file",
			opts:    &CreateOptions{CustomFieldsFile: fieldsPath},
			wantLen: 1,
		},
		{
			name: "both sources",
			opts: &CreateOptions{
				CustomFieldsJSON: `[{"id":"field","value":"demo"}]`,
				CustomFieldsFile: fieldsPath,
			},
			wantErr: "cannot use both --custom-fields-json and --custom-fields-file",
		},
		{
			name:    "invalid json",
			opts:    &CreateOptions{CustomFieldsJSON: `{"id":"field"}`},
			wantErr: "invalid custom fields JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getCustomFields(tt.opts)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("getCustomFields() error = %v, want substring %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("getCustomFields() error = %v", err)
			}
			if len(got) != tt.wantLen {
				t.Fatalf("len(customFields) = %d, want %d", len(got), tt.wantLen)
			}
		})
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func issueResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Status:     http.StatusText(status),
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}
