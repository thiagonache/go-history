package history

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
)

type recorder struct {
	Ctx        context.Context
	File       io.WriteCloser
	Path       string
	Permission os.FileMode
	Signals    []os.Signal
	Stdin      io.Reader
	Stdout     io.Writer
	Stop       context.CancelFunc
}

// NewRecorder instantiate a new Recorder object and returns a pointer to it.
func NewRecorder() (*recorder, error) {
	return &recorder{
		Path:       "history.log",
		Permission: 0664,
		Signals:    []os.Signal{os.Interrupt},
		Stdout:     os.Stdout,
		Stdin:      os.Stdin,
	}, nil
}

// EnsureHistoryFileOpen ensures the recorder log file is opened before writing
// to it. It does allow the user to overwrite the default file path.
func (r *recorder) EnsureHistoryFileOpen() error {
	if r.File != nil {
		return nil
	}
	history, err := os.OpenFile(r.Path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, r.Permission)
	if err != nil {
		return err
	}
	r.File = history
	return nil
}

// Session reads user input which should be a shell unix command, executes it
// and record the commands and its outputs
func (r *recorder) Session() {
	err := r.EnsureHistoryFileOpen()
	if err != nil {
		fmt.Fprint(r.Stdout, err)
	}
	tee := io.MultiWriter(r.Stdout, r.File)
	for {
		fmt.Fprint(tee, "$ ")
		reader := bufio.NewReader(r.Stdin)
		input, err := reader.ReadString('\n')
		if err == io.EOF {
			r.Stop()
			break
		}
		if err != nil {
			fmt.Fprint(tee, err)
		}
		input = input[:len(input)-1]
		if input == "exit" || input == "quit" {
			fmt.Fprintln(tee, input)
			r.Stop()
		}
		fmt.Fprintln(r.File, input)
		// err = r.Execute(input)
		// if err == ? { // need to handle Disk full
		// 	r.Stop()
		// }
		r.Execute(input)
	}
}

// Execute receives the command to run, executes it and implements
// io.MultiWriter to write to Recorder Stdout and Recorder file.
func (r *recorder) Execute(command string) error {
	tee := io.MultiWriter(r.Stdout, r.File)
	entrypoint := strings.Split(command, " ")[0]
	args := strings.Split(command, " ")[1:]
	cmd := exec.Command(entrypoint, args...)
	cmd.Stderr = tee
	cmd.Stdout = tee
	cmd.Stdin = r.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// ListenSignals creates the parent context in order to handle OS Signals
func (r *recorder) ListenSignals() {
	r.Ctx, r.Stop = signal.NotifyContext(context.Background(), r.Signals...)
}

// Shutdown implements a graceful shutdown for the package by displaying the
// path of the file with the data recorded and make sure the file descriptor is
// closed.
func (r recorder) Shutdown() {
	fmt.Fprintf(r.Stdout, "\rSee recorded data at %s\n", r.Path)
	err := r.File.Close()
	if err != nil {
		log.Fatal(err)
	}
}
