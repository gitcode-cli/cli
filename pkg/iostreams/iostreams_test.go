package iostreams

import "testing"

func TestSplitCommandLine(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []string
		wantErr bool
	}{
		{
			name:  "simple",
			input: "less -R",
			want:  []string{"less", "-R"},
		},
		{
			name:  "double quotes",
			input: `delta --syntax-theme="GitHub Dark"`,
			want:  []string{"delta", "--syntax-theme=GitHub Dark"},
		},
		{
			name:  "single quotes",
			input: `pager --title 'hello world'`,
			want:  []string{"pager", "--title", "hello world"},
		},
		{
			name:  "escaped spaces",
			input: `pager /tmp/with\ spaces/file.txt`,
			want:  []string{"pager", "/tmp/with spaces/file.txt"},
		},
		{
			name:    "unterminated quote",
			input:   `pager "oops`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := splitCommandLine(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("splitCommandLine() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(got) != len(tt.want) {
				t.Fatalf("splitCommandLine() = %#v, want %#v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("splitCommandLine() = %#v, want %#v", got, tt.want)
				}
			}
		})
	}
}
