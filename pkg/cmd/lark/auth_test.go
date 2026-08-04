package lark

import (
	"strings"
	"testing"

	"gitcode.com/gitcode-cli/cli/pkg/iostreams"
)

func TestAuthStatusRun_FriendlyIdentity(t *testing.T) {
	io, _, out, _ := iostreams.Test()
	opts := &authStatusOptions{
		LarkRun: stubRun(`{"ok":true,"identity":"user"}`, "", 0),
		IO:      io,
	}
	if err := authStatusRun(opts); err != nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(out.String(), "logged in as user") {
		t.Errorf("out = %q, want identity summary", out.String())
	}
}

func TestAuthStatusRun_JSONPassthrough(t *testing.T) {
	payload := `{"ok":true,"identity":"user","identities":{"user":{"status":"ready"}}}`
	io, _, out, _ := iostreams.Test()
	opts := &authStatusOptions{
		LarkRun: stubRun(payload, "", 0),
		JSON:    true,
		IO:      io,
	}
	if err := authStatusRun(opts); err != nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(out.String(), `"identities"`) {
		t.Errorf("out = %q, want passthrough of identities", out.String())
	}
}

func TestAuthStatusRun_StdoutEnvelope(t *testing.T) {
	// lark-cli writes auth status JSON to stdout; identity should still resolve.
	io, _, out, _ := iostreams.Test()
	opts := &authStatusOptions{
		LarkRun: stubRun(`{"ok":true,"identity":"bot"}`, "", 0),
		IO:      io,
	}
	if err := authStatusRun(opts); err != nil {
		t.Fatalf("err = %v", err)
	}
	if !strings.Contains(out.String(), "logged in as bot") {
		t.Errorf("out = %q, want bot identity", out.String())
	}
}

func TestAuthStatusRun_NonZeroExit(t *testing.T) {
	io, _, _, _ := iostreams.Test()
	opts := &authStatusOptions{
		LarkRun: stubRun("", `{"ok":false,"error":{"message":"not logged in"}}`, 1),
		IO:      io,
	}
	err := authStatusRun(opts)
	if err == nil {
		t.Fatal("err = nil, want error")
	}
	if !strings.Contains(err.Error(), "not logged in") {
		t.Errorf("err = %v, want not logged in", err)
	}
}
