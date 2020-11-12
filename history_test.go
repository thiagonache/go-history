package history_test

import (
	"bytes"
	"history"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExecuteAndRecordCommand(t *testing.T) {
	entrypoint := "echo"
	args := []string{"testing"}
	var output bytes.Buffer
	err := history.ExecuteAndRecordCommand(&bytes.Buffer{}, &output, entrypoint, args...)
	if err != nil {
		t.Fatal(err)
	}
	wantOutput := "testing\n"
	if !cmp.Equal(wantOutput, output.String()) {
		t.Error(cmp.Diff(wantOutput, output.String()))
	}
}

func TestRecordSession(t *testing.T) {
	var fakeOutput bytes.Buffer
	var historyBuf bytes.Buffer
	var fakeInput bytes.Buffer
	_, err := fakeInput.WriteString("echo testing\nexit\n")
	if err != nil {
		t.Fatalf("cannot write string to the buffer: %v", err)
	}
	err = history.RecordSession(&fakeInput, &fakeOutput, &historyBuf)
	if err != nil {
		t.Fatal(err)
	}
	wantHistory := "$ echo testing\ntesting\n$ exit\n"
	if !cmp.Equal(wantHistory, historyBuf.String()) {
		t.Error(cmp.Diff(wantHistory, historyBuf.String()))
	}
	wantOutput := "$ testing\n$ "
	if !cmp.Equal(wantOutput, fakeOutput.String()) {
		t.Error(cmp.Diff(wantOutput, fakeOutput.String()))
	}
}
