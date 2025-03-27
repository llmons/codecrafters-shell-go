package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/codecrafters-io/shell-starter-go/app/excuter"
	"github.com/codecrafters-io/shell-starter-go/app/scanner"
	"github.com/codecrafters-io/shell-starter-go/app/util"
)

func main() {
	for {
		// print prompt and get command
		fmt.Print("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		scanner := scanner.NewScanner(command)
		tokens := scanner.ScanTokens()

		if len(tokens) == 0 {
			continue
		}

		if excuter, err := excuter.NewExecutor(tokens); err == nil {
			result := excuter.Execute()
			if result != "" {
				fmt.Println(result)
			}
		} else {
			util.Error(err.Error())
		}
	}
}
