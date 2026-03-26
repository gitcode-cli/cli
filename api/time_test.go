package api

import (
	"encoding/json"
	"testing"
	"time"
)

func TestFlexibleTime_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		wantZero bool
		wantYear int
	}{
		{
			name:     "RFC3339 with timezone",
			json:     `"2026-03-26T16:03:07+08:00"`,
			wantZero: false,
			wantYear: 2026,
		},
		{
			name:     "RFC3339 UTC",
			json:     `"2026-03-26T08:03:07Z"`,
			wantZero: false,
			wantYear: 2026,
		},
		{
			name:     "ISO 8601 without timezone",
			json:     `"2026-03-26T16:03:07"`,
			wantZero: false,
			wantYear: 2026,
		},
		{
			name:     "Date only",
			json:     `"2026-04-30"`,
			wantZero: false,
			wantYear: 2026,
		},
		{
			name:     "Common datetime format",
			json:     `"2026-03-26 16:03:07"`,
			wantZero: false,
			wantYear: 2026,
		},
		{
			name:     "Null value",
			json:     `null`,
			wantZero: true,
			wantYear: 0,
		},
		{
			name:     "Empty string",
			json:     `""`,
			wantZero: true,
			wantYear: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ft FlexibleTime
			err := json.Unmarshal([]byte(tt.json), &ft)
			if err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}

			if ft.IsZero() != tt.wantZero {
				t.Errorf("IsZero() = %v, want %v", ft.IsZero(), tt.wantZero)
			}

			if !tt.wantZero && ft.Year() != tt.wantYear {
				t.Errorf("Year() = %v, want %v", ft.Year(), tt.wantYear)
			}
		})
	}
}

func TestFlexibleTime_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		time    FlexibleTime
		wantNull bool
	}{
		{
			name:    "Zero time",
			time:    FlexibleTime{},
			wantNull: true,
		},
		{
			name:    "Valid time",
			time:    FlexibleTime{Time: time.Date(2026, 3, 26, 16, 3, 7, 0, time.UTC)},
			wantNull: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.time)
			if err != nil {
				t.Fatalf("Failed to marshal: %v", err)
			}

			if tt.wantNull && string(data) != "null" {
				t.Errorf("Expected null, got %s", string(data))
			}

			if !tt.wantNull && string(data) == "null" {
				t.Errorf("Expected non-null value, got null")
			}
		})
	}
}

func TestIssue_Unmarshal_FlexibleTime(t *testing.T) {
	jsonResp := `{
		"id": 123456,
		"number": "1",
		"title": "Test Issue",
		"body": "Test body",
		"state": "open",
		"html_url": "https://gitcode.com/test/test/issues/1",
		"user": {
			"id": 1,
			"login": "testuser"
		},
		"created_at": "2026-03-26",
		"updated_at": "2026-04-30T15:00:00Z"
	}`

	var issue Issue
	err := json.Unmarshal([]byte(jsonResp), &issue)
	if err != nil {
		t.Fatalf("Failed to unmarshal Issue: %v", err)
	}

	if issue.Number != "1" {
		t.Errorf("Expected Number '1', got '%s'", issue.Number)
	}

	if issue.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if issue.CreatedAt.Year() != 2026 {
		t.Errorf("Expected CreatedAt year 2026, got %d", issue.CreatedAt.Year())
	}

	if issue.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}

	if issue.UpdatedAt.Year() != 2026 {
		t.Errorf("Expected UpdatedAt year 2026, got %d", issue.UpdatedAt.Year())
	}
}

func TestPullRequest_Unmarshal_FlexibleTime(t *testing.T) {
	jsonResp := `{
		"id": 123456,
		"number": 1,
		"title": "Test PR",
		"body": "Test body",
		"state": "open",
		"html_url": "https://gitcode.com/test/test/pull/1",
		"user": {
			"id": 1,
			"login": "testuser"
		},
		"created_at": "2026-03-26T10:00:00Z",
		"updated_at": "2026-04-30"
	}`

	var pr PullRequest
	err := json.Unmarshal([]byte(jsonResp), &pr)
	if err != nil {
		t.Fatalf("Failed to unmarshal PullRequest: %v", err)
	}

	if pr.Number != 1 {
		t.Errorf("Expected Number 1, got %d", pr.Number)
	}

	if pr.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}

	if pr.UpdatedAt.IsZero() {
		t.Error("UpdatedAt should not be zero")
	}
}

func TestCommitAuthor_Unmarshal_FlexibleTime(t *testing.T) {
	jsonResp := `{
		"name": "Test User",
		"email": "test@example.com",
		"date": "2026-03-26"
	}`

	var author CommitAuthor
	err := json.Unmarshal([]byte(jsonResp), &author)
	if err != nil {
		t.Fatalf("Failed to unmarshal CommitAuthor: %v", err)
	}

	if author.Name != "Test User" {
		t.Errorf("Expected Name 'Test User', got '%s'", author.Name)
	}

	if author.Date.IsZero() {
		t.Error("Date should not be zero")
	}

	if author.Date.Year() != 2026 {
		t.Errorf("Expected Date year 2026, got %d", author.Date.Year())
	}
}