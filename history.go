package history

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// LogFile is the name where the recorded data will be stored
const LogFile = "history.log"

// LogPerm is the file permission in unix format
const LogPerm = 0644

// RecordSession takes an io.Reader and an io.Writer, reads the input up to the new line,
// call ExecuteAndRecordCommand function and writes the output into the io.Writer.
// An error is returned if it happens otherwise nil.
func RecordSession(r io.Reader, output io.Writer, history io.Writer) error {
	var err error
	// The CLI send nil to tell this function that it should create a log file.
	if history == nil {
		fmt.Println("Welcome to history")
		fmt.Printf("See recorded data at %s\n", LogFile)
		history, err = os.OpenFile(LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, LogPerm)
		if err != nil {
			return err
		}
	}
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
		err = ExecuteAndRecordCommand(r, tee, entrypoint, args...)
		if err != nil {
			return err
		}
	}
}

// ExecuteAndRecordCommand takes an io.Writer (stdin or bytes.buffer), an
// entrypoint and args to call RunCommand function. An error is returned if
// found otherwise nil.
func ExecuteAndRecordCommand(r io.Reader, output io.Writer, entrypoint string, args ...string) error {
	// ioErr stores any error when writing to the io.Writer
	var ioErr error
	// When the command return an error we store and print the error. Otherwise,
	// we store and print the command output.
	cmd := exec.Command(entrypoint, args...)
	cmd.Stderr = output
	cmd.Stdout = output
	cmd.Stdin = r
	err := cmd.Run()
	if err != nil {
		_, ioErr = fmt.Fprintln(output, err.Error())
		fmt.Println(err.Error())
	}
	if ioErr != nil {
		return ioErr
	}

	return nil
}
