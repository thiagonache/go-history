package history_test

import (
	"bytes"
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

	command := "echo testing"
	want := "testing\n"
	r, err := history.NewRecorder(
		history.WithLogPath("/tmp/history.log"),
	)
	if err != nil {
		t.Fatal(err)
	}
	output := &bytes.Buffer{}
	r.Stdout = output
	r.Stderr = output
	historyBuf := &testWriteCloser{}
	r.File = historyBuf
	r.Execute(command)
	got := historyBuf.b.String()
	if !cmp.Equal(want, got) {
		t.Error(cmp.Diff(want, got))
	}
}

func TestErrorsCmdNotExist(t *testing.T) {
	t.Parallel()

	r, err := history.NewRecorder(
		history.WithLogPath("/tmp/history.log"),
	)
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

	r, err := history.NewRecorder(
		history.WithLogPath("/tmp/history.log"),
	)
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
	r.WaitForExit()
}
