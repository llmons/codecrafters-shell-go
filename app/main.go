package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var builtinCommands map[string]func([]string) error
var HOME string

func init() {
	builtinCommands = map[string]func([]string) error{
		"exit": handleExit,
		"echo": handleEcho,
		"type": handleType,
		"pwd":  handlePwd,
		"cd":   handleCd,
	}
	HOME = os.Getenv("HOME")
}

func parseCommand(command string) (string, []string, error) {
	// get tokens
	tokens := []string{}
	command = strings.TrimSpace(command)
	for {
		start := strings.Index(command, "'")
		if start == -1 {
			tokens = append(tokens, strings.Fields(command)...)
			break
		}
		tokens = append(tokens, strings.Fields(command[:start])...)
		command = command[start+1:]
		end := strings.Index(command, "'")
		tokens = append(tokens, command[:end])
		command = command[end+1:]
	}

	if len(tokens) == 0 {
		return "", nil, fmt.Errorf("empty command")
	}
	majorCommand := strings.ToLower(tokens[0])
	argv := tokens[1:]
	return majorCommand, argv, nil
}

func handleExit(argv []string) error {
	if len(argv) == 0 {
		os.Exit(0)
	}
	if code, err := strconv.Atoi(argv[0]); err == nil {
		os.Exit(code)
		return nil
	} else {
		return err
	}
}

func handleEcho(argv []string) error {
	if len(argv) == 0 {
		fmt.Println()
		return nil
	}
	fmt.Println(strings.Join(argv, " "))
	return nil
}

func handleType(argv []string) error {
	if len(argv) == 0 {
		return nil
	}
	if _, ok := builtinCommands[argv[0]]; ok {
		fmt.Println(argv[0] + " is a shell builtin")
		return nil
	}
	if path, err := exec.LookPath(argv[0]); err == nil {
		fmt.Println(argv[0] + " is " + path)
		return nil
	} else {
		fmt.Println(argv[0] + ": not found")
		return nil
	}
}

func handlePwd(argv []string) error {
	if len(argv) != 0 {
		return fmt.Errorf("pwd: too many arguments")
	}
	if dir, err := os.Getwd(); err == nil {
		fmt.Println(dir)
		return nil
	}
	return fmt.Errorf("pwd: error")
}

func handleCd(argv []string) error {
	if len(argv) > 1 {
		return fmt.Errorf("cd: too many arguments")
	}
	if len(argv) == 0 || argv[0] == "~" {
		if err := os.Chdir(HOME); err != nil {
			fmt.Println("cd: " + HOME + ": No such file or directory")
		}
	} else if err := os.Chdir(argv[0]); err != nil {
		fmt.Println("cd: " + argv[0] + ": No such file or directory")
	}
	return nil
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
		if command == "\n" {
			continue
		}
		majorCommand, argv, err := parseCommand(command)
		if err != nil {
			fmt.Println(err)
			continue
		}

		// handle builtin commands
		if handler, ok := builtinCommands[majorCommand]; ok {
			handler(argv)
			continue
		}

		// check if major command exists
		if _, err := exec.LookPath(majorCommand); err != nil {
			fmt.Println(majorCommand + ": command not found")
			continue
		}

		// execute common command
		cmd := exec.Command(majorCommand, argv...)
		stdout, err := cmd.Output()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Print(string(stdout))
		}
	}
}
