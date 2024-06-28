package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/cobra"

	"github.com/reddtsai/goREPL/pkg/validation"
)

type Repl struct {
	mu         sync.Mutex
	users      map[string]bool
	validation validation.IValidation
	rootCmd    *cobra.Command
}

func New() *Repl {
	repl := &Repl{
		users:      make(map[string]bool),
		validation: validation.NewValidation(),
	}
	repl.rootCmd = &cobra.Command{
		Use:   "repl",
		Short: "virtual file system (REPL)",
		Run:   repl.RootCmdRunner,
	}

	return repl
}
func (r *Repl) Execute() error {
	return r.rootCmd.Execute()
}

func (r *Repl) RootCmdRunner(cmd *cobra.Command, args []string) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("# ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "exit" {
			fmt.Println("Goodbye!")
			break
		}
		args := strings.Fields(line)
		if len(args) > 0 {
			cmd.SetArgs(args)
			cmd.Execute()
		}
	}
}

func (r *Repl) AddRegisterCmd() {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register a user",
		Run:   r.RegisterRunner,
	}

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) RegisterRunner(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Fprintln(os.Stderr, "invalid command")
		return
	}

	// case insensitive
	userName := strings.ToLower(args[0])
	// input validation
	validate := r.validation.ValidateUserName(userName)
	if !validate.IsSuccess {
		fmt.Fprintln(os.Stderr, validate.Message)
		return
	}
	_, exist := r.users[userName]
	if exist {
		msg := fmt.Sprintf("The [%s] has already existed", userName)
		fmt.Fprintln(os.Stderr, msg)
		return
	}
	// save
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[userName] = true

	fmt.Printf("Add %s successfully\n", userName)
}
