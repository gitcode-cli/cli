package api

import (
	"encoding/json"
	"testing"
)

func TestFlexibleNumber_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "string number", input: `"123"`, want: "123"},
		{name: "integer number", input: `123`, want: "123"},
		{name: "zero", input: `0`, want: "0"},
		{name: "large number", input: `99999`, want: "99999"},
		{name: "null", input: `null`, want: ""},
		{name: "string zero", input: `"0"`, want: "0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var fn FlexibleNumber
			err := json.Unmarshal([]byte(tt.input), &fn)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON(%s) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if string(fn) != tt.want {
				t.Errorf("UnmarshalJSON(%s) = %q, want %q", tt.input, string(fn), tt.want)
			}
		})
	}
}

func TestFlexibleNumber_UnmarshalJSON_Struct(t *testing.T) {
	// Full struct decode: ensures FlexibleNumber works inside Issue
	t.Run("string number in struct", func(t *testing.T) {
		j := `{"number":"42","title":"Bug"}`
		var issue struct {
			Number FlexibleNumber `json:"number"`
			Title  string         `json:"title"`
		}
		if err := json.Unmarshal([]byte(j), &issue); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if string(issue.Number) != "42" {
			t.Errorf("Number = %q, want 42", issue.Number)
		}
	})

	t.Run("integer number in struct", func(t *testing.T) {
		j := `{"number":42,"title":"Bug"}`
		var issue struct {
			Number FlexibleNumber `json:"number"`
			Title  string         `json:"title"`
		}
		if err := json.Unmarshal([]byte(j), &issue); err != nil {
			t.Fatalf("Unmarshal error: %v", err)
		}
		if string(issue.Number) != "42" {
			t.Errorf("Number = %q, want 42", issue.Number)
		}
	})
}

func TestFlexibleNumber_MarshalJSON(t *testing.T) {
	fn := FlexibleNumber("42")
	data, err := json.Marshal(fn)
	if err != nil {
		t.Fatalf("Marshal error: %v", err)
	}
	if string(data) != `"42"` {
		t.Errorf("MarshalJSON() = %s, want \"42\"", data)
	}
}

func TestFlexibleNumber_String(t *testing.T) {
	fn := FlexibleNumber("42")
	if fn.String() != "42" {
		t.Errorf("String() = %q, want 42", fn.String())
	}
}
