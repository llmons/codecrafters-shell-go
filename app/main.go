package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint
var builtinCommands = []string{"exit", "echo", "type"}

func parseCommand(command string) (string, []string, error) {
	if command[len(command)-1] == '\n' {
		command = command[:len(command)-1]
	}
	argvs := strings.Split(command, " ")
	if len(argvs) == 0 {
		return "", nil, fmt.Errorf("empty command")
	}
	majorCommand := argvs[0]
	argvs = argvs[1:]
	return majorCommand, argvs, nil
}

func execBuiltinCommand(majorCommand string, argvs []string) error {
	switch majorCommand {
	// exit command
	case "exit":
		if len(argvs) == 0 {
			os.Exit(0)
		}
		if argvs[0] == "0" {
			os.Exit(0)
		}
		os.Exit(1)

	// echo command
	case "echo":
		if len(argvs) == 0 {
			fmt.Println()
			return nil
		}
		fmt.Println(strings.Join(argvs, " "))
		return nil

	// type command
	case "type":
		if len(argvs) == 0 {
			return nil
		}
		if slices.Contains(builtinCommands, argvs[0]) {
			fmt.Println(argvs[0] + " is a shell builtin")
			return nil
		}
		if path, err := exec.LookPath(argvs[0]); err == nil {
			fmt.Println(argvs[0] + " is " + path)
			return nil
		}
		fmt.Println(argvs[0] + ": not found")
		return nil
	}

	return fmt.Errorf("execBuiltinCommand: unknown command")
}

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		command = strings.Trim(command, " ")
		if command == "\n" {
			continue
		}

		majorCommand, argvs, err := parseCommand(command)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if slices.Contains(builtinCommands, majorCommand) {
			if err := execBuiltinCommand(majorCommand, argvs); err != nil {
				fmt.Println(err)
			}
			continue
		}

		if _, err := exec.LookPath(majorCommand); err != nil {
			fmt.Println(majorCommand + ": command not found")
			continue
		}

		cmd := exec.Command(majorCommand, argvs...)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Print(string(stdout))
		}
	}
}
