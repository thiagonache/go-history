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
	"path/filepath"
	"strings"
)

// Recorder store object data for the package
type Recorder struct {
	Context    context.Context
	File       io.WriteCloser
	path       string
	permission os.FileMode
	Stderr     io.Writer
	Stdin      io.Reader
	Stdout     io.Writer
	Stop       context.CancelFunc
}

// NewRecorder instantiate a new Recorder object and returns a pointer to it.
func NewRecorder() (*Recorder, error) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	return &Recorder{
		Context:    ctx,
		path:       "history.log",
		permission: 0664,
		Stderr:     os.Stderr,
		Stdin:      os.Stdin,
		Stdout:     os.Stdout,
		Stop:       stop,
	}, nil
}

func isDirectory(path string) error {
	basedir := filepath.Dir(path)
	info, err := os.Stat(basedir)
	if os.IsNotExist(err) {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is a file and must be a directory", basedir)
	}
	return nil
}

// EnsureHistoryFileOpen ensures the recorder log file is opened before writing
// to it. It does allow the user to overwrite the default file path.
func (r *Recorder) EnsureHistoryFileOpen() error {
	// Check if file descriptor is already opened.
	if r.File != nil {
		return nil
	}
	err := isDirectory(r.path)
	if err != nil {
		return err
	}
	history, err := os.OpenFile(r.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, r.permission)
	if err != nil {
		return err
	}
	r.File = history
	return nil
}

// Session reads user input which should be a shell unix command, executes it
// and record the commands and its outputs
func (r *Recorder) Session() {
	err := r.EnsureHistoryFileOpen()
	if err != nil {
		r.Stop()
		fmt.Fprintln(r.Stderr, err)
		os.Exit(1)
	}
	tee := io.MultiWriter(r.Stdout, r.File)
	teeErr := io.MultiWriter(r.Stderr, r.File)
	for {
		fmt.Fprint(tee, "$ ")
		reader := bufio.NewReader(r.Stdin)
		input, err := reader.ReadString('\n')
		if err == io.EOF {
			r.Stop()
			break
		}
		if err != nil {
			fmt.Fprint(teeErr, err)
		}
		input = input[:len(input)-1]
		if input == "exit" || input == "quit" {
			fmt.Fprintln(tee, input)
			r.Stop()
		}
		fmt.Fprintln(r.File, input)
		err = r.Execute(input)
		if err != nil {
			fmt.Fprintln(teeErr, err)
		}
	}
}

// Execute receives the command to run, executes it and implements
// io.MultiWriter to write to Recorder Stdout and Recorder file.
func (r *Recorder) Execute(command string) error {
	tee := io.MultiWriter(r.Stdout, r.File)
	teeErr := io.MultiWriter(r.Stderr, r.File)
	entrypoint := strings.Split(command, " ")[0]
	args := strings.Split(command, " ")[1:]
	cmd := exec.Command(entrypoint, args...)
	cmd.Stderr = teeErr
	cmd.Stdout = tee
	cmd.Stdin = r.Stdin
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

// SetPath takes a string and set history log file path
func (r *Recorder) SetPath(path string) {
	r.path = path
}

// SetPermission takes a string and set history log file path
func (r *Recorder) SetPermission(perm os.FileMode) {
	r.permission = perm
}

// Shutdown implements a graceful shutdown for the package by displaying the
// path of the file with the data recorded and make sure the file descriptor is
// closed.
func (r Recorder) Shutdown() {
	fmt.Fprintf(r.Stdout, "\rSee recorded data at %s\n", r.path)
	err := r.File.Close()
	if err != nil {
		log.Fatal(err)
	}
}
