package main

import (
	"log"

	"github.com/reddtsai/goREPL/cmd"
)

func main() {
	repl := cmd.New()
	repl.AddRegisterCmd()
	err := repl.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
