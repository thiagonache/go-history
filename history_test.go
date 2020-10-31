package history_test

import (
	"bytes"
	"history"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRunCommand(t *testing.T) {
	wantEntrypoint := "command"
	wantParam1 := "param1"
	cmdLine := "command param1 paramX"
	gotEntrypoint := strings.Split(cmdLine, " ")[0]
	gotParam1 := strings.Split(cmdLine, " ")[1]

	if wantEntrypoint != gotEntrypoint {
		t.Errorf("want %s, got %s", wantEntrypoint, gotEntrypoint)
	}
	if wantParam1 != gotParam1 {
		t.Errorf("want %s, got %s", wantParam1, gotParam1)
	}
}

func TestExecuteAndRecordCommand(t *testing.T) {
	command := "echo testing"
	var got bytes.Buffer
	err := history.ExecuteAndRecordCommand(command, &got)
	if err != nil {
		t.Fatal(err)
	}
	want := "echo testing\ntesting\n"

	if !cmp.Equal(want, got.String()) {
		t.Error(cmp.Diff(want, got.String()))
	}
}

func TestRun(t *testing.T) {
	var got bytes.Buffer
	var stdin bytes.Buffer
	// For some reason via tests it does return error io.EOF
	// when reading the reader even having \n at the end
	_, err := stdin.WriteString("echo hello\n")
	if err != nil {
		t.Fatalf("cannot write string to the buffer: %e", err)
	}
	err = history.Run(&stdin, &got)
	if err != nil {
		t.Fatal(err)
	}
	want := "abc"

	if want != got.String() {
		t.Fatalf("want %q and got %q", want, got.String())
	}
}
