package api

import (
	"fmt"
	"net/http"
)

// User represents a GitCode user
type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
	HTMLURL   string `json:"html_url"`
	CreatedAt string `json:"created_at"`
}

// CurrentUser fetches the current authenticated user
func CurrentUser(client *Client) (*User, error) {
	var user User
	err := client.Get("/user", &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUser fetches a user by username
func GetUser(client *Client, username string) (*User, error) {
	var user User
	err := client.Get(fmt.Sprintf("/users/%s", username), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// VerifyToken verifies that a token is valid by fetching the current user
func VerifyToken(httpClient *http.Client, host, token string) (*User, error) {
	client := NewClient(httpClient, host, token)
	return CurrentUser(client)
}

// ClientFromToken creates a client with the given token
func ClientFromToken(token string) *Client {
	return NewClient(DefaultHTTPClient(), DefaultHost, token)
}

// ClientFromTokenAndHost creates a client with the given token and host
func ClientFromTokenAndHost(token, host string) *Client {
	return NewClient(DefaultHTTPClient(), host, token)
}