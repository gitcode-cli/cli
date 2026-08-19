package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

// SSHKey represents an SSH public key registered with GitCode.
type SSHKey struct {
	ID        int64  `json:"id"`
	Title     string `json:"title"`
	Key       string `json:"key"`
	CreatedAt string `json:"created_at"`
	URL       string `json:"url"`
}

// CreateSSHKeyOptions contains the fields required to register an SSH key.
type CreateSSHKeyOptions struct {
	Title string `json:"title"`
	Key   string `json:"key"`
}

// ListSSHKeys lists SSH public keys for the authenticated user.
func ListSSHKeys(client *Client, page, perPage int) ([]SSHKey, error) {
	path := "/user/keys" + newQueryBuilder().SetInt("page", page).SetInt("per_page", perPage).String()
	var keys sshKeysResponse
	if err := client.Get(path, &keys); err != nil {
		return nil, err
	}
	return keys.Keys, nil
}

// CreateSSHKey registers an SSH public key for the authenticated user.
func CreateSSHKey(client *Client, opts *CreateSSHKeyOptions) (*SSHKey, error) {
	var key SSHKey
	if err := client.Post("/user/keys", opts, &key); err != nil {
		return nil, err
	}
	return &key, nil
}

// GetSSHKey fetches an SSH public key by ID.
func GetSSHKey(client *Client, id int64) (*SSHKey, error) {
	var response sshKeyResponse
	if err := client.Get("/user/keys/"+strconv.FormatInt(id, 10), &response); err != nil {
		return nil, err
	}
	if response.Key == nil {
		return nil, fmt.Errorf("ssh key %d: empty response", id)
	}
	return response.Key, nil
}

// DeleteSSHKey deletes an SSH public key by ID.
func DeleteSSHKey(client *Client, id int64) error {
	return client.Delete("/user/keys/" + strconv.FormatInt(id, 10))
}

// The published OpenAPI schema describes GET /user/keys/{id} as an array,
// while its examples contain a single object. Accept both server shapes.
type sshKeyResponse struct{ Key *SSHKey }
type sshKeysResponse struct{ Keys []SSHKey }

func (r *sshKeysResponse) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}
	if data[0] == '[' {
		return json.Unmarshal(data, &r.Keys)
	}
	var key SSHKey
	if err := json.Unmarshal(data, &key); err != nil {
		return err
	}
	r.Keys = []SSHKey{key}
	return nil
}

func (r *sshKeyResponse) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 || bytes.Equal(data, []byte("null")) {
		return nil
	}
	if data[0] == '[' {
		var keys []SSHKey
		if err := json.Unmarshal(data, &keys); err != nil {
			return err
		}
		if len(keys) > 0 {
			r.Key = &keys[0]
		}
		return nil
	}
	var key SSHKey
	if err := json.Unmarshal(data, &key); err != nil {
		return err
	}
	r.Key = &key
	return nil
}
