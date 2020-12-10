package history

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// LogPerm is the file permission in unix format
const LogPerm = 0644

type Recorder struct {
	Path string
	File io.Writer
	Stdout io.Writer
	Stdin io.Reader
}

func NewRecorder() (*Recorder, error) {
	return &Recorder{
		Path: "history.log",
		Stdout: os.Stdout,
		Stdin: os.Stdin,
	}, nil
}

func (r *Recorder) EnsureHistoryFileOpen() error {
	if r.File != nil {
		return nil
	}
	history, err := os.OpenFile(r.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogPerm)
	if err != nil {
		return err
	}
	r.File = history
	return nil
}

// RecordSession takes an io.Reader and an io.Writer, reads the input up to the new line,
// call ExecuteAndRecordCommand function and writes the output into the io.Writer.
// An error is returned if it happens otherwise nil.
func (r *Recorder) Session() error {
	var err error
	fmt.Fprintln(r.Stdout, "Welcome to history")
	fmt.Fprintf(r.Stdout, "See recorded data at %s\n", r.LogFile)
	tee := io.MultiWriter(output, history)
	for {
		fmt.Fprint(tee, "$ ")
		reader := bufio.NewReader(r)
		input, err := reader.ReadString('\n')
		// When control+d is pressed we get EOF which should be handled gracefully
		if err == io.EOF {
			return nil
		}
		if err != nil {
			// %w preserve error type
			return fmt.Errorf("error reading the input: %w", err)
		}
		input = input[:len(input)-1]
		if input == "exit" || input == "quit" {
			return nil
		}
		fmt.Fprintln(history, input)
		entrypoint := strings.Split(input, " ")[0]
		args := strings.Split(input, " ")[1:]
		err = Execute(entrypoint, args...)
		if err != nil {
			return err
		}
	}
}

// Execute takes an io.Writer (stdin or bytes.buffer), an
// entrypoint and args to call RunCommand function. An error is returned if
// found otherwise nil.
func (r *Recorder) Execute(entrypoint string, args ...string) error {
	// When the command return an error we store and print the error.
	// Otherwise, we store and print the command output.
	cmd := exec.Command(entrypoint, args...)
	cmd.Stderr = r.Stdout
	cmd.Stdout = r.Stdout
	cmd.Stdin = r.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
