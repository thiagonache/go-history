package history_test

import (
	"bytes"
	"fmt"
	"testing"

	history "github.com/thiagonache/go-history"

	"github.com/google/go-cmp/cmp"
)

type testWriteCloser struct {
	bytes.Buffer
}

func (wrc *testWriteCloser) Close() error { return nil }

func TestExecute(t *testing.T) {
	testCases := []struct {
		command     string
		desc        string
		errExpected bool
		output      bytes.Buffer
		want        string
	}{
		{
			desc:    "Simple Echo",
			command: "echo testing",
			want:    "testing\n",
		},
		{
			desc:        "Non-existing command",
			command:     "abc",
			want:        "exec: \"abc\": executable file not found in $PATH",
			errExpected: true,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			r, err := history.NewRecorder()
			if err != nil {
				t.Fatal(err)
			}
			// EnsureHistoryFileOpen is required here because it happens on the
			// Session function in order to follow the designed flow.
			err = r.EnsureHistoryFileOpen()
			if err != nil {
				fmt.Fprint(r.Stdout, err)
			}
			r.Stdout = &tC.output
			err = r.Execute(tC.command)
			errFound := err != nil
			output := tC.output.String()
			if tC.errExpected {
				if tC.errExpected != errFound {
					t.Fatalf("unexpected error")
				}
				output = err.Error()
			}

			if !cmp.Equal(tC.want, output) {
				t.Error(cmp.Diff(tC.want, output))
			}
			r.Shutdown()
		})
	}
}
func TestSession(t *testing.T) {
	var fakeOutput bytes.Buffer
	var historyBuf testWriteCloser
	var fakeInput bytes.Buffer

	_, err := fakeInput.WriteString("echo testing\nexit\n")
	if err != nil {
		t.Fatalf("cannot write string to the buffer: %v", err)
	}

	r, err := history.NewRecorder()
	if err != nil {
		t.Fatal(err)
	}

	r.Stdin = &fakeInput
	r.Stdout = &fakeOutput
	r.File = &historyBuf
	r.Session()

	wantHistory := "$ echo testing\ntesting\n$ "
	if !cmp.Equal(wantHistory, historyBuf.String()) {
		t.Error(cmp.Diff(wantHistory, historyBuf.String()))
	}
	wantOutput := "$ testing\n$ "
	if !cmp.Equal(wantOutput, fakeOutput.String()) {
		t.Error(cmp.Diff(wantOutput, fakeOutput.String()))
	}
	r.Shutdown()
}
