package history

import (
	"bufio"
	"fmt"
	"io"
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
	input, err := reader.ReadString('\n')
	// When control+d is pressed we get EOF which should be handled gracefully
	if err == io.EOF {
		return err
	}
	if err != nil {
		return fmt.Errorf("error reading the input: %e", err)
	}
	input = input[:len(input)-1]
	if input == "exit" || input == "quit" {
		return io.EOF
	}
	entrypoint := strings.Split(input, " ")[0]
	args := strings.Split(input, " ")[1:]
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
		fmt.Fprintln(w, err.Error())
		fmt.Println(err.Error())
	} else {
		// output already have new line at the end
		fmt.Fprint(w, output)
		fmt.Print(output)
	}

	return nil
}
