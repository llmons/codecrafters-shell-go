package main

import (
	"bufio"
	"fmt"
	"github.com/codecrafters-io/shell-starter-go/app/excuter"
	"github.com/codecrafters-io/shell-starter-go/app/scanner"
	"github.com/codecrafters-io/shell-starter-go/app/util"
	"os"
)

func main() {
	for {
		// print prompt and get command
		fmt.Println("$ ")
		command, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		scanner := scanner.NewScanner(command)
		tokens := scanner.ScanTokens()

		if excuter, err := excuter.NewExecutor(tokens); err == nil {
			result := excuter.Execute()
			fmt.Println(result)
		} else {
			util.Error(err.Error())
		}
	}
}
