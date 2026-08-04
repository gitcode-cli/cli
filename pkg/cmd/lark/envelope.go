package lark

import "encoding/json"

// larkEnvelope models the JSON contract documented by lark-cli: success and
// error responses are distinct envelopes. Success goes to stdout with ok=true;
// errors go to stderr with ok=false and an error object. gc parses these to
// normalize output for its own --json contract.
type larkEnvelope struct {
	OK       bool            `json:"ok"`
	Identity string          `json:"identity,omitempty"`
	Data     json.RawMessage `json:"data,omitempty"`
	Error    *larkError      `json:"error,omitempty"`
	Meta     json.RawMessage `json:"meta,omitempty"`
}

// larkError is the error sub-object inside a lark-cli error envelope.
type larkError struct {
	Type    string `json:"type,omitempty"`
	Subtype string `json:"subtype,omitempty"`
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
	Hint    string `json:"hint,omitempty"`
}

// parseEnvelope decodes a lark-cli JSON payload. It tolerates empty/invalid
// payloads by returning a zero envelope plus the decode error.
func parseEnvelope(raw []byte) (larkEnvelope, error) {
	var env larkEnvelope
	if len(raw) == 0 {
		return env, nil
	}
	err := json.Unmarshal(raw, &env)
	return env, err
}
