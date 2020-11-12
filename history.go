package history

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// RunCommand takes an entrypoint and arguments, execute the command, and
// returns the command output and an error if it happens.
func RunCommand(entrypoint string, args []string) (string, error) {
	cmd := exec.Command(entrypoint, args...)
	var stdout bytes.Buffer
	if entrypoint == "vim" {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = &stdout
	}
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return stdout.String(), nil
}

// Run takes an io.Reader and an io.Writer, reads the input up to the new line,
// call ExecuteAndRecordCommand function and writes the output into the io.Writer.
// An error is returned if it happens otherwise nil.
func Run(r io.Reader, w io.Writer) error {
	for {
		fmt.Fprint(os.Stdout, "$ ")
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
		entrypoint := strings.Split(input, " ")[0]
		args := strings.Split(input, " ")[1:]
		err = ExecuteAndRecordCommand(w, entrypoint, args...)
		if err != nil {
			return err
		}
	}
}

// ExecuteAndRecordCommand takes an io.Writer (stdin or bytes.buffer), an
// entrypoint and args to call RunCommand function. An error is returned if
// found otherwise nil.
func ExecuteAndRecordCommand(w io.Writer, entrypoint string, args ...string) error {
	fmt.Fprintf(w, entrypoint)
	for _, arg := range args {
		fmt.Fprintf(w, " "+arg)
	}
	fmt.Fprint(w, "\n")

	// ioErr stores any error when writing to the io.Writer
	var ioErr error
	// When the command return an error we store and print the error. Otherwise,
	// we store and print the command output.
	output, err := RunCommand(entrypoint, args)
	if err != nil {
		_, ioErr = fmt.Fprintln(w, err.Error())
		fmt.Println(err.Error())
	} else {
		// output already have a new line at the end.
		_, ioErr = fmt.Fprint(w, output)
		fmt.Print(output)
	}
	if ioErr != nil {
		return ioErr
	}

	return nil
}
