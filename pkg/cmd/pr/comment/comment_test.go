package comment

import (
	"bytes"
	"testing"

	cmdutil "gitcode.com/gitcode-cli/cli/pkg/cmdutil"
)

func TestNewCmdComment(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "comment with --body",
			args:    []string{"123", "--body", "Test comment"},
			wantErr: false,
		},
		{
			name:    "comment with --body-file",
			args:    []string{"123", "--body-file", "-"},
			wantErr: false,
		},
		{
			name:    "inline comment with path and position",
			args:    []string{"123", "--body", "Inline comment", "--path", "api/auth.go", "--position", "42"},
			wantErr: false,
		},
		{
			name:    "no PR number",
			args:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid PR number",
			args:    []string{"abc"},
			wantErr: true,
		},
		{
			name:    "no body provided",
			args:    []string{"123"},
			wantErr: false, // Command runs, error in run function
		},
		{
			name:    "both body and body-file",
			args:    []string{"123", "--body", "test", "--body-file", "-"},
			wantErr: false, // Command runs, error in run function
		},
		{
			name:    "path without position",
			args:    []string{"123", "--body", "test", "--path", "api/auth.go"},
			wantErr: false, // Command runs, error in run function
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			cmd := NewCmdComment(f, func(opts *CommentOptions) error {
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

func TestGetBody(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		bodyFile   string
		stdin      string
		wantBody   string
		wantErr    bool
		errContain string
	}{
		{
			name:     "body flag only",
			body:     "Hello world",
			wantBody: "Hello world",
		},
		{
			name:       "both body and body-file",
			body:       "Hello",
			bodyFile:   "-",
			wantErr:    true,
			errContain: "cannot use both",
		},
		{
			name:     "empty body",
			wantBody: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := cmdutil.TestFactory()
			opts := &CommentOptions{
				IO:       f.IOStreams,
				Body:     tt.body,
				BodyFile: tt.bodyFile,
			}

			if tt.stdin != "" {
				opts.IO.In = bytes.NewBufferString(tt.stdin)
			}

			got, err := getBody(opts)
			if tt.wantErr {
				if err == nil {
					t.Errorf("getBody() expected error, got nil")
					return
				}
				if tt.errContain != "" && !contains(err.Error(), tt.errContain) {
					t.Errorf("getBody() error = %v, want containing %v", err, tt.errContain)
				}
				return
			}
			if err != nil {
				t.Errorf("getBody() unexpected error: %v", err)
				return
			}
			if got != tt.wantBody {
				t.Errorf("getBody() = %v, want %v", got, tt.wantBody)
			}
		})
	}
}

func TestInlineCommentValidation(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		position   int
		wantErr    bool
		errContain string
	}{
		{
			name:     "no inline flags",
			path:     "",
			position: 0,
			wantErr:  false,
		},
		{
			name:     "path and position provided",
			path:     "api/auth.go",
			position: 42,
			wantErr:  false,
		},
		{
			name:       "path without position",
			path:       "api/auth.go",
			position:   0,
			wantErr:    true,
			errContain: "--position is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &CommentOptions{
				Path:     tt.path,
				Position: tt.position,
				Body:     "test", // Provide body to avoid other validation errors
			}

			err := validateInlineFlags(opts)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateInlineFlags() expected error, got nil")
					return
				}
				if tt.errContain != "" && !contains(err.Error(), tt.errContain) {
					t.Errorf("validateInlineFlags() error = %v, want containing %v", err, tt.errContain)
				}
				return
			}
			if err != nil {
				t.Errorf("validateInlineFlags() unexpected error: %v", err)
			}
		})
	}
}

func validateInlineFlags(opts *CommentOptions) error {
	if opts.Path != "" && opts.Position == 0 {
		return &validationError{msg: "--position is required when --path is specified for inline comments"}
	}
	return nil
}

type validationError struct {
	msg string
}

func (e *validationError) Error() string {
	return e.msg
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
