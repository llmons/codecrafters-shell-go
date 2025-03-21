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

func main() {
	builtinCommands := []string{"exit", "echo", "type"}

	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		majorCommand, _, err := parseCommand(command)
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch majorCommand {
		// exit command
		case "exit":
			if command == "exit 0\n" {
				os.Exit(0)
			}

		// echo command
		case "echo":
			if len(command) > 5 && command[:4] == "echo" {
				fmt.Println(command[5 : len(command)-1])
				continue
			}

		// type command
		case "type":
			if len(command) > 5 && command[:4] == "type" {
				if target := command[5 : len(command)-1]; slices.Contains(builtinCommands, target) {
					fmt.Println(target + " is a shell builtin")
				} else if path, err := exec.LookPath(target); err == nil {
					fmt.Println(target + " is " + path)
				} else {
					fmt.Println(target + ": not found")
				}
				continue
			}

		default:
			fmt.Println(command[:len(command)-1] + ": command not found")
		}
	}
}
