package excuter

import (
	"errors"
	"github.com/codecrafters-io/shell-starter-go/app/scanner"
	"os"
	"os/exec"
	"slices"
	"strconv"
	"strings"
)

var builtinCommands = []string{
	"cd", "echo", "exit", "pwd", "type",
}

const HOME = "HOME"

type Executor struct {
	command string
	argv    []string
	result  string
}

func NewExecutor(tokens []scanner.Token) (*Executor, error) {
	command := tokens[0].Lexeme
	if _, err := exec.LookPath(command); !isBuiltin(command) && err != nil {
		return nil, errors.New("command not found")
	}

	argv := make([]string, len(tokens)-1)
	for i, token := range tokens[1:] {
		argv[i] = token.Lexeme
	}
	return &Executor{
		command: tokens[0].Lexeme,
		argv:    argv,
	}, nil
}

func (e *Executor) Execute() string {
	switch e.command {
	case "cd":
		e.cd()
	case "echo":
		e.echo()
	case "exit":
		e.exit()
	case "pwd":
		e.pwd()
	case "type":
		e._type()
	default:
		cmd := exec.Command(e.command, e.argv...)
		if stdout, err := cmd.Output(); err != nil {
			e.result = err.Error()
		} else {
			e.result = string(stdout)
		}
	}

	return e.result
}

func (e *Executor) cd() {
	if len(e.argv) > 1 {
		e.result = "cd: too many arguments"
		return
	}
	if len(e.argv) == 0 || e.argv[0] == "~" {
		if err := os.Chdir(HOME); err != nil {
			e.result = "cd: " + HOME + ": No such file or directory"
		}
	}
	if err := os.Chdir(e.argv[0]); err != nil {
		e.result = "cd: " + e.argv[0] + ": No such file or directory"
	}
}

func (e *Executor) echo() {
	if len(e.argv) == 0 {
		e.result = ""
		return
	}
	e.result = strings.Join(e.argv, " ")
}

func (e *Executor) exit() {
	if len(e.argv) == 0 {
		os.Exit(0)
	}
	if code, err := strconv.Atoi(e.argv[0]); err == nil {
		os.Exit(code)
	}
}

func (e *Executor) pwd() {
	if len(e.argv) != 0 {
		e.result = "pwd: too many arguments"
		return
	}
	if dir, err := os.Getwd(); err == nil {
		e.result = dir
		return
	}
	e.result = "pwd: error"
}

func (e *Executor) _type() {
	if len(e.argv) == 0 {
		e.result = ""
		return
	}
	if isBuiltin(e.argv[0]) {
		e.result = e.argv[0] + " is a shell builtin"
		return
	}
	if path, err := exec.LookPath(e.argv[0]); err == nil {
		e.result = e.argv[0] + " is " + path
		return
	}
	e.result = e.argv[0] + ": not found"
}

func isBuiltin(command string) bool {
	return slices.Contains(builtinCommands, command)
}
