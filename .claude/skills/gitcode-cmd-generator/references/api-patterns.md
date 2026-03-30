# GitCode API Patterns

> This file is a Claude reference layer for command generation.
> API usage patterns must remain consistent with `spec/`, current code in `api/`, and actual command behavior documented in `docs/COMMANDS.md`.

This reference documents common patterns for implementing API operations in gitcode-cli.

## API Client Usage

### Basic Setup

```go
import (
	"net/http"
	"github.com/gitcode-com/gitcode-cli/api"
)

// In command run function:
httpClient, err := opts.HttpClient()
if err != nil {
	return fmt.Errorf("failed to create HTTP client: %w", err)
}

client := api.NewClientFromHTTP(httpClient)
token := getEnvToken()
if token == "" {
	return fmt.Errorf("not authenticated. Run: gc auth login")
}
client.SetToken(token, "environment")
```

### Token Retrieval

```go
func getEnvToken() string {
	if token := os.Getenv("GC_TOKEN"); token != "" {
		return token
	}
	return os.Getenv("GITCODE_TOKEN")
}
```

## Common API Operations

### GET Request

```go
// Get single resource
var result ResourceType
err := client.Get("/repos/"+owner+"/"+repo+"/resource/"+id, &result)

// Get list of resources
var results []ResourceType
err := client.Get("/repos/"+owner+"/"+repo+"/resources", &results)
```

### POST Request

```go
opts := &CreateOptions{
	Field1: "value1",
	Field2: "value2",
}
var result ResourceType
err := client.Post("/repos/"+owner+"/"+repo+"/resources", opts, &result)
```

### PATCH Request

```go
opts := &UpdateOptions{
	Field1: "new value",
}
var result ResourceType
err := client.Patch("/repos/"+owner+"/"+repo+"/resources/"+id, opts, &result)
```

### DELETE Request

```go
err := client.Delete("/repos/" + owner + "/" + repo + "/resources/" + id)
```

## Data Structures

### User

```go
type User struct {
	ID        interface{} `json:"id"`
	Login     string      `json:"login"`
	Name      string      `json:"name"`
	Email     string      `json:"email"`
	AvatarURL string      `json:"avatar_url"`
	HTMLURL   string      `json:"html_url"`
}
```

### Repository

```go
type Repository struct {
	ID          interface{} `json:"id"`
	Name        string      `json:"name"`
	FullName    string      `json:"full_name"`
	Description string      `json:"description"`
	Private     bool        `json:"private"`
	HTMLURL     string      `json:"html_url"`
	CloneURL    string      `json:"clone_url"`
	SSHURL      string      `json:"ssh_url"`
	Owner       *User       `json:"owner"`
}
```

### Pagination

```go
type ListOptions struct {
	PerPage int `url:"per_page,omitempty"`
	Page    int `url:"page,omitempty"`
}

// Default pagination
const DefaultPerPage = 30
const MaxPerPage = 100
```

## API Endpoints Reference

| Resource | Endpoint | Methods |
|----------|----------|---------|
| User | `/user` | GET |
| Repositories | `/user/repos`, `/repos/{owner}/{repo}` | GET, POST, DELETE |
| Issues | `/repos/{owner}/{repo}/issues` | GET, POST |
| Issue | `/repos/{owner}/{repo}/issues/{number}` | GET, PATCH |
| Comments | `/repos/{owner}/{repo}/issues/{number}/comments` | GET, POST |
| Pull Requests | `/repos/{owner}/{repo}/pulls` | GET, POST |
| PR | `/repos/{owner}/{repo}/pulls/{number}` | GET, PATCH |
| Labels | `/repos/{owner}/{repo}/labels` | GET, POST, DELETE |
| Milestones | `/repos/{owner}/{repo}/milestones` | GET, POST, DELETE |

## Error Handling

```go
// API errors
var apiErr *api.APIError
if errors.As(err, &apiErr) {
	if apiErr.StatusCode == 404 {
		return fmt.Errorf("resource not found")
	}
	if apiErr.StatusCode == 403 {
		return fmt.Errorf("permission denied")
	}
	return fmt.Errorf("API error: %s", apiErr.Message)
}

// Generic errors
if err != nil {
	return fmt.Errorf("operation failed: %w", err)
}
```

## Helper Functions

### Parse Repository

```go
func parseRepo(repo string) (string, string, error) {
	if repo == "" {
		// TODO: get from current git repo
		return "", "", fmt.Errorf("no repository specified. Use -R owner/repo")
	}

	parts := strings.Split(repo, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid repository format: %s", repo)
	}
	return parts[0], parts[1], nil
}
```

### Integer Conversion

```go
// itoa converts int to string (defined in api package)
func itoa(i int) string {
	return strconv.Itoa(i)
}

// atoi converts string to int with error handling
func atoi(s string) (int, error) {
	return strconv.Atoi(s)
}
```
