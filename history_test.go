package history_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
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
		history.WithLogPath("/tmp/history-test-execute.log"),
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
	r.Shutdown()
}

func TestErrorsCmdNotExist(t *testing.T) {
	t.Parallel()

	r, err := history.NewRecorder(
		history.WithLogPath("/tmp/history-cmd-not-exist.log"),
	)
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
	r.Shutdown()
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
		history.WithLogPath("/tmp/history-session.log"),
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

func TestWithLogPath(t *testing.T) {
	t.Parallel()
	const charset = "abcdefghijklmnopqrstuvxzABCDEFGHIJKLMNOPQRSTUVXZ"
	b := make([]byte, 5)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	s := string(b)
	want := fmt.Sprintf("%s.log", s)
	historyFile := fmt.Sprintf("/tmp/%s.log", s)
	r, err := history.NewRecorder(
		history.WithLogPath(historyFile),
	)
	if err != nil {
		t.Fatal(err)
	}
	fd, err := os.Stat(historyFile)
	if err != nil {
		t.Fatal(err)
	}
	got := fd.Name()
	if !cmp.Equal(want, got) {
		t.Errorf(cmp.Diff(want, got))
	}
	r.Shutdown()
}

func TestWithLogPermission(t *testing.T) {
	t.Parallel()

	r, err := history.NewRecorder(
		history.WithLogPermission(0600),
		history.WithLogPath("/tmp/history-log-permission.log"),
	)
	if err != nil {
		t.Fatal(err)
	}
	var want os.FileMode = 0600
	fd, err := os.Stat("/tmp/history-log-permission.log")
	if err != nil {
		t.Fatal(err)
	}
	got := fd.Mode().Perm()

	if want != got {
		t.Errorf("want %d, got %d", want, got)
	}
	r.Shutdown()
}
