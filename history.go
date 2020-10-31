package history

import (
	"bufio"
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
	// How can I print the $ only on CLI?
	// Should I move the for loop to the CLI?
	for {
		reader := bufio.NewReader(r)
		text, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading the input: %e", err)
		}
		text = text[:len(text)-1]
		if text == "exit" || text == "quit" {
			break
		}
		err = ExecuteAndRecordCommand(text, w)
		if err != nil {
			return err
		}
	}

	return nil
}

func ExecuteAndRecordCommand(cmd string, w io.Writer) error {
	fmt.Fprintln(w, cmd)
	entrypoint := strings.Split(cmd, " ")[0]
	args := strings.Split(cmd, " ")[1:]
	output, err := RunCommand(entrypoint, args)
	fmt.Fprint(w, output)
	if err == io.ErrUnexpectedEOF { // need to confirm io error for out of space
		return err
	}
	if err != nil {
		fmt.Fprintf(w, err.Error())
		fmt.Println(err.Error())
	} else {
		fmt.Println(output)
	}

	return nil
}
