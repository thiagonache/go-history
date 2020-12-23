package history_test

import (
	"bytes"
	"fmt"
	"history"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type TestWriteCloser struct {
	bytes.Buffer
}

func (wrc *TestWriteCloser) Close() error { return nil }

func TestExecuteAndRecordCommand(t *testing.T) {
	command := "echo testing"
	var output bytes.Buffer
	r, err := history.NewRecorder()
	if err != nil {
		t.Fatal(err)
	}
	err = r.EnsureHistoryFileOpen()
	if err != nil {
		fmt.Fprint(r.Stdout, err)
	}
	r.Stdout = &output
	err = r.Execute(command)
	if err != nil {
		t.Fatal(err)
	}
	wantOutput := "testing\n"
	if !cmp.Equal(wantOutput, output.String()) {
		t.Error(cmp.Diff(wantOutput, output.String()))
	}
	r.Shutdown()
}

func TestRecordSession(t *testing.T) {
	var fakeOutput bytes.Buffer
	var historyBuf TestWriteCloser
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
	r.ListenSignals()
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
