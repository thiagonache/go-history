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
	if err != nil {
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
	fmt.Fprint(w, "\n")
	output, err := RunCommand(entrypoint, args)
	// When the OS returns error we need to exit
	if err == io.ErrClosedPipe || err == io.ErrShortWrite || err == io.ErrUnexpectedEOF {
		return err
	}
	// When the command return an error we store and print the error. Otherwise,
	// we store and print the command output
	if err != nil {
		fmt.Fprint(w, err.Error())
		fmt.Println(err.Error())
	} else {
		fmt.Fprint(w, output)
		// output already have new line at the end
		fmt.Printf(output)
	}

	return nil
}
