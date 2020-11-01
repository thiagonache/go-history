package history

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

func RunCommand(entrypoint string, args []string) (string, error) {
	output, err := exec.Command(entrypoint, args...).Output()
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func Run(r io.Reader, w io.Writer) error {
	reader := bufio.NewReader(r)
	text, err := reader.ReadString('\n')
	if err != nil && err != io.EOF {
		return fmt.Errorf("error reading the input: %e", err)
	}
	text = text[:len(text)-1]
	if text == "exit" {
		return errors.New("exit")
	}
	entrypoint := strings.Split(string(text), " ")[0]
	args := strings.Split(string(text), " ")[1:]
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
		fmt.Println(output)
	}

	return nil
}
