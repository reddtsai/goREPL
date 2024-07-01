package main

import (
	"log"

	"github.com/reddtsai/goREPL/cmd"
)

func main() {
	repl := cmd.New()
	repl.AddRegisterCmd()     // 1
	repl.AddCreateFolderCmd() // 2
	repl.AddDeleteFolderCmd() // 3
	repl.AddListFoldersCmd()  // 4
	repl.AddRenameFolderCmd() // 5
	repl.AddCreateFileCmd()   // 6
	repl.AddDeleteFileCmd()   // 7
	repl.AddListFilesCmd()    // 8

	err := repl.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
