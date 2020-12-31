package history_test

import (
	"bytes"
	"fmt"
	"testing"

	history "github.com/thiagonache/go-history"

	"github.com/google/go-cmp/cmp"
)

type testWriteCloser struct {
	b bytes.Buffer
}

func (wrc *testWriteCloser) Write(data []byte) (int, error) {
	return wrc.b.Write(data)
}

func (wrc *testWriteCloser) Len() int {
	return wrc.b.Len()
}

func (wrc *testWriteCloser) Close() error { return nil }

func TestExecute(t *testing.T) {
	t.Parallel()

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
			command:     "doesntexist",
			want:        "BOGUS",
			errExpected: true,
		},
	}
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
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output := &bytes.Buffer{}
			r.Stdout = output
			err = r.Execute(tC.command)
			errFound := err != nil
			got := output.String()
			if tC.errExpected != errFound {
				t.Fatalf("unexpected error status %v", err)
			}
			if !tC.errExpected && !cmp.Equal(tC.want, got) {
				t.Error(cmp.Diff(tC.want, got))
			}
		})
	}
	r.Shutdown()
}

func TestErrorsToHistory(t *testing.T) {
	r, err := history.NewRecorder()
	if err != nil {
		t.Fatal(err)
	}
	err = r.EnsureHistoryFileOpen()
	if err != nil {
		t.Fatal(err)
	}
	fakeStdErr := &bytes.Buffer{}
	r.Stderr = fakeStdErr
	historyBuf := &testWriteCloser{}
	r.File = historyBuf
	err = r.Execute("doesntexist")
	if err == nil {
		t.Fatal(err)
	}
	if fakeStdErr.Len() == 0 {
		t.Error("want something written to stderr, got nothing")
	}
	if historyBuf.Len() == 0 {
		t.Error("want something written to history file, got nothing")
	}
}

func TestSession(t *testing.T) {
	t.Parallel()

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
	if !cmp.Equal(wantHistory, historyBuf.b.String()) {
		t.Error(cmp.Diff(wantHistory, historyBuf.b.String()))
	}
	wantOutput := "$ testing\n$ "
	if !cmp.Equal(wantOutput, fakeOutput.String()) {
		t.Error(cmp.Diff(wantOutput, fakeOutput.String()))
	}
	r.Shutdown()
}
