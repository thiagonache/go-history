package history_test

import (
	"bytes"
	"history"
	"strings"
	"testing"

	"github.com/bitfield/script"
	"github.com/google/go-cmp/cmp"
)

func TestWriteHistoryFile(t *testing.T) {
	var want string = "rm /tmp/some_file"
	var filePath = "/tmp/file"

	strCmd := []string{"ls /tmp", "rm /tmp/some_file"}
	err := history.WriteFile(filePath, strCmd)
	if err != nil {
		t.Errorf("cannot write to the file: %v", err)
	}
	got, err := script.File(filePath).Match(want).String()
	if err != nil {
		t.Errorf("cannot match %s on file %s: %w", want, filePath, err)
	}
	if want != strings.TrimSuffix(got, "\n") {
		t.Errorf("want %q, got %q", want, got)
	}
}

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

func TestInputRead(t *testing.T) {
	want := "Hello world!\n"
	got, err := history.ReadInputFrom(bytes.NewBufferString(want))
	if err != nil {
		t.Errorf("error reading input: %v", err)
	}
	if want != got {
		t.Errorf("want %s, got %s.", want, got)
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
