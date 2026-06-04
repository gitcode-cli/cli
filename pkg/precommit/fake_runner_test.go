package precommit

import "strings"

// fakeRunner is a deterministic CommandRunner for tests.
type fakeRunner struct {
	look      map[string]bool     // executables present on PATH
	responses map[string]fakeResp // keyed by "name arg1 arg2"
	calls     []string            // recorded "name arg1 arg2" invocations
}

type fakeResp struct {
	out    string // stdout
	stderr string // stderr; combined into Run output but excluded from RunStdout
	err    error
}

func newFakeRunner() *fakeRunner {
	return &fakeRunner{
		look:      map[string]bool{},
		responses: map[string]fakeResp{},
	}
}

func key(name string, args ...string) string {
	if len(args) == 0 {
		return name
	}
	return name + " " + strings.Join(args, " ")
}

func (f *fakeRunner) Look(name string) bool { return f.look[name] }

func (f *fakeRunner) Run(_ string, name string, args ...string) (string, error) {
	k := key(name, args...)
	f.calls = append(f.calls, k)
	r := f.responses[k]
	return r.out + r.stderr, r.err
}

func (f *fakeRunner) RunStdout(_ string, name string, args ...string) (string, error) {
	k := key(name, args...)
	f.calls = append(f.calls, k)
	r := f.responses[k]
	return r.out, r.err
}

func (f *fakeRunner) called(name string, args ...string) bool {
	want := key(name, args...)
	for _, c := range f.calls {
		if c == want {
			return true
		}
	}
	return false
}
