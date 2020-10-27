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

func Run() error {
	var cmdHistory []string
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
		cmdHistory = append(cmdHistory, fmt.Sprintf("command: %s", text))
		output, err := RunCommand(strings.Split(text, " ")[0], strings.Split(text, " ")[1:])
		if err != nil {
			fmt.Printf("[ERROR] cannot run command %s: %v\n", text, err)
			cmdHistory = append(cmdHistory, fmt.Sprintf("error: %s", err.Error()))
			continue
		}
		fmt.Printf(output)
		cmdHistory = append(cmdHistory, fmt.Sprintf("output: %s", output))
	}
	WriteFile(".history", cmdHistory)
	return nil
}
