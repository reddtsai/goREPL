package main

import (
	"log"

	"github.com/reddtsai/goREPL/cmd"
)

func main() {
	repl := cmd.New()
	repl.AddRegisterCmd()
	repl.AddCreateFolderCmd()
	repl.AddDeleteFolderCmd()
	repl.AddListFoldersCmd()
	repl.AddRenameFolderCmd()
	repl.AddCreateFileCmd()
	repl.AddDeleteFileCmd()
	repl.AddListFilesCmd()

	err := repl.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}
