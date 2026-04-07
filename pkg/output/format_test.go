package output

import "testing"

func TestParseFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Format
		wantErr bool
	}{
		{name: "default", input: "", want: FormatSimple},
		{name: "json", input: "json", want: FormatJSON},
		{name: "simple", input: "simple", want: FormatSimple},
		{name: "table", input: "table", want: FormatTable},
		{name: "invalid", input: "yaml", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("ParseFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestParseTimeFormat(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    TimeFormat
		wantErr bool
	}{
		{name: "default", input: "", want: TimeFormatAbsolute},
		{name: "absolute", input: "absolute", want: TimeFormatAbsolute},
		{name: "relative", input: "relative", want: TimeFormatRelative},
		{name: "invalid", input: "iso", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseTimeFormat(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseTimeFormat() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Fatalf("ParseTimeFormat() = %q, want %q", got, tt.want)
			}
		})
	}
}
