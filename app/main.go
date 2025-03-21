package main

import (
	"bufio"
	"fmt"
	"os"
)

// Ensures gofmt doesn't remove the "fmt" import in stage 1 (feel free to remove this!)
var _ = fmt.Fprint

func main() {
	for {
		fmt.Fprint(os.Stdout, "$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// exit command
		if command == "exit 0\n" {
			os.Exit(0)
		}

		// echo command
		if len(command) > 5 && command[:4] == "echo" {
			fmt.Println(command[5 : len(command)-1])
			continue
		}

		fmt.Println(command[:len(command)-1] + ": command not found")
	}
}
