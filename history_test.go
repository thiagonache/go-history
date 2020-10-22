package history_test

import (
	"bytes"
	"history"
	"os/exec"
	"testing"
)

func TestWriteHistoryFile(t *testing.T) {
	var want string = "rm /tmp/some_file\n"
	var grepFor = "/tmp/some_file"
	var filePath = "/tmp/file"

	strCmd := []string{"ls /tmp", "rm /tmp/some_file"}
	err := history.WriteFile(filePath, strCmd)
	if err != nil {
		t.Errorf("cannot write to the file: %v", err)
	}
	output, err := exec.Command("grep", grepFor, "/tmp/file").Output()
	if err != nil {
		t.Errorf("cannot run grep %s %s: %v", grepFor, "/tmp/file", err)
	}
	got := string(output)
	if want != got {
		t.Errorf("want %q, got %q", want, got)
	}
}

func TestRunCommand(t *testing.T) {
	want := "Hello world!\n"
	cmd := "echo"
	args := []string{"Hello", "world!"}
	got, err := history.RunCommand(cmd, args)
	if err != nil {
		t.Errorf("cannot run command %s due to %v", cmd, err)
	}
	if want != got {
		t.Errorf("want %s, got %s", want, got)
	}
}

func TestInputRead(t *testing.T) {
	want := "Hello world!\n"
	got, err := history.ReadInput(bytes.NewBufferString(want))
	if err != nil {
		t.Errorf("error reading input: %v", err)
	}
	if want != got {
		t.Errorf("want %s, got %s.", want, got)
	}
}
