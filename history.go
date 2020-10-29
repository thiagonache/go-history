package history

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func WriteFile(filePath string, lines []string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("cannot write to file %s, error %v", filePath, err)
	}
	for _, line := range lines {
		f.WriteString(fmt.Sprintf("%s%s", line, "\n"))
	}
	defer f.Close()

	return nil
}

func ReadInputFrom(r io.Reader) (string, error) {
	reader := bufio.NewReader(r)
	text, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("cannot read input: %v", err)
	}
	return text, nil
}

func RunCommand(entrypoint string, args []string) (string, error) {
	output, err := exec.Command(entrypoint, args...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func Run(w io.Writer) error {
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Printf("$ ")
		text, err := reader.ReadString('\n')
		text = text[:len(text)-1]
		if err != nil {
			return fmt.Errorf("error reading the input")
		}
		if text == "exit" || text == "quit" {
			break
		}

		ExecuteAndRecordCommand(text, w)
	}
	return nil
}

func ExecuteAndRecordCommand(cmd string, w io.Writer) error {
	fmt.Fprintln(w, cmd)
	entrypoint := strings.Split(cmd, " ")[0]
	args := strings.Split(cmd, " ")[1:]
	output, err := RunCommand(entrypoint, args)
	fmt.Fprint(w, output)
	if err != nil {
		return err
	}
	return nil
}
