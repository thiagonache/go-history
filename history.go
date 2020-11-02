package history

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// RunCommand takes an entrypoint and arguments, execute the command, and
// returns the command output and error
func RunCommand(entrypoint string, args []string) (string, error) {
	output, err := exec.Command(entrypoint, args...).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

// Run takes an io.Reader and an io.Writer, reads the input up to the new line
// and call ExecuteAndRecordCommand function
func Run(r io.Reader, w io.Writer) error {
	reader := bufio.NewReader(r)
	text, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading the input: %e", err)
	}
	text = text[:len(text)-1]
	if text == "exit" || text == "quit" {
		os.Exit(0)
	}
	entrypoint := strings.Split(text, " ")[0]
	args := strings.Split(text, " ")[1:]
	err = ExecuteAndRecordCommand(w, entrypoint, args...)
	if err != nil {
		return err
	}

	return nil
}

// ExecuteAndRecordCommand takes an io.Writer (stdin or bytes.buffer), an
// entrypoint and args to compose the full command.
func ExecuteAndRecordCommand(w io.Writer, entrypoint string, args ...string) error {
	fmt.Fprintf(w, entrypoint)
	for _, arg := range args {
		fmt.Fprintf(w, " "+arg)
	}
	fmt.Fprintf(w, "\n")
	output, err := RunCommand(entrypoint, args)
	fmt.Fprint(w, output)
	// I need help to figure this out
	// if err == outOfSpace {
	// 	return err
	// }
	if err != nil {
		fmt.Fprintf(w, "%s\n", err.Error())
		fmt.Println(err.Error())
	} else {
		fmt.Printf(output)
	}

	return nil
}
