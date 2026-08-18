package delete

import (
	"io"
	"net/http"
	"strings"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
	"gitcode.com/gitcode-cli/cli/pkg/testutil"
)

func TestNewCmdDeleteParsesID(t *testing.T) {
	cmd := NewCmdDelete(cmdutil.TestFactory(), func(opts *DeleteOptions) error {
		if opts.ID != 42 || !opts.Yes || !opts.JSON {
			t.Fatalf("options = %+v", opts)
		}
		return nil
	})
	cmd.SetArgs([]string{"42", "--yes", "--json"})
	if err := cmd.Execute(); err != nil {
		t.Fatal(err)
	}
}

func TestDeleteRunRequiresConfirmation(t *testing.T) {
	f := cmdutil.TestFactory()
	err := deleteRun(&DeleteOptions{IO: f.IOStreams, ID: 42})
	if err == nil || !strings.Contains(err.Error(), "confirmation required") {
		t.Fatalf("error = %v", err)
	}
}

func TestDeleteRunDryRunDoesNotCallAPI(t *testing.T) {
	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	f.IOStreams.Out = out
	err := deleteRun(&DeleteOptions{IO: f.IOStreams, ID: 42, DryRun: true, HttpClient: func() (*http.Client, error) {
		t.Fatal("HttpClient must not be called during dry-run")
		return nil, nil
	}})
	if err != nil || !strings.Contains(out.String(), "Would delete SSH key 42") {
		t.Fatalf("output = %q, error = %v", out, err)
	}
}

func TestDeleteRunJSON(t *testing.T) {
	t.Setenv("GC_TOKEN", "test-token")
	f := cmdutil.TestFactory()
	out := &strings.Builder{}
	f.IOStreams.Out = out
	opts := &DeleteOptions{IO: f.IOStreams, ID: 42, Yes: true, JSON: true, HttpClient: func() (*http.Client, error) {
		return &http.Client{Transport: testutil.NewRoundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodDelete || req.URL.Path != "/api/v5/user/keys/42" {
				t.Fatalf("request = %s %s", req.Method, req.URL.Path)
			}
			return &http.Response{StatusCode: http.StatusOK, Status: "OK", Header: make(http.Header), Body: io.NopCloser(strings.NewReader("{}"))}, nil
		})}, nil
	}}
	if err := deleteRun(opts); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "\"action\": \"deleted\"") {
		t.Fatalf("output = %s", out)
	}
}
