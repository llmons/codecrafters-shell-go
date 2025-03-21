package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
)

var builtinCommands = []string{"exit", "echo", "type", "pwd", "cd"}
var HOME = os.Getenv("HOME")

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
	argc := len(argvs)
	switch majorCommand {
	// exit command
	case "exit":
		if argc == 0 || argvs[0] == "0" {
			os.Exit(0)
		}
		os.Exit(1)

	// echo command
	case "echo":
		if argc == 0 {
			fmt.Println()
			return nil
		}
		fmt.Println(strings.Join(argvs, " "))
		return nil

	// type command
	case "type":
		if argc == 0 {
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

	// pwd command
	case "pwd":
		if argc != 0 {
			return fmt.Errorf("pwd: too many arguments")
		}
		if dir, err := os.Getwd(); err == nil {
			fmt.Println(dir)
		}
		return nil

	// cd command
	case "cd":
		if argc == 0 {
			return nil
		}
		if argc > 1 {
			return fmt.Errorf("cd: too many arguments")
		}
		if argvs[0] == "~" {
			if err := os.Chdir(HOME); err != nil {
				return fmt.Errorf("cd: %s: No such file or directory", HOME)
			}
			return nil
		}
		if err := os.Chdir(argvs[0]); err != nil {
			return fmt.Errorf("cd: %s: No such file or directory", argvs[0])
		}
		return nil
	}

	return fmt.Errorf("execBuiltinCommand: unknown command")
}

func main() {
	for {
		// print prompt and get command
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// handle command string
		command = strings.Trim(command, " ")
		if command == "\n" {
			continue
		}
		majorCommand, argvs, err := parseCommand(command)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// builtin commands
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
