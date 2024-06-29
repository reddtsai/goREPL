package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"

	"github.com/reddtsai/goREPL/pkg/storage"
)

type Repl struct {
	storage           storage.IStorage
	rootCmd           *cobra.Command
	folderSortName    string
	folderSortCreated string
}

// New returns a new Repl
func New() *Repl {
	repl := &Repl{
		storage: storage.NewVirtualFileSysStorage(),
	}
	repl.rootCmd = &cobra.Command{
		Use:     "repl",
		Version: "1.0.0",
		Short:   "virtual file system (REPL)",
		Run:     repl.RootCmdRunner,
	}

	return repl
}

// Execute runs the REPL
func (r *Repl) Execute() error {
	return r.rootCmd.Execute()
}

func (r *Repl) RootCmdRunner(cmd *cobra.Command, args []string) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("======== Virtual File System 1.0.0 ========")
	fmt.Println("Please Enter Your Command")

	for {
		fmt.Print("# ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		switch line {
		case "exit":
			fmt.Println("Goodbye!")
			return
		case "help":
			fmt.Println("> register [username]")
			fmt.Println("> create-folder [username] [foldername] [description]?")
		default:
			args := r.SplitArgs(line)
			if len(args) > 0 {
				foundCmd, _, err := cmd.Find(args)
				if err != nil || foundCmd == r.rootCmd {
					fmt.Fprintln(os.Stderr, "Error: Invalid command.")
					continue
				}
				cmd.SetArgs(args)
				_ = cmd.Execute()
			}
		}
	}
}

func (r *Repl) SplitArgs(line string) []string {
	var args []string
	var currentArg string
	var inQuote bool

	for _, r := range line {
		switch {
		case unicode.IsSpace(r) && !inQuote:
			if currentArg != "" {
				args = append(args, currentArg)
				currentArg = ""
			}
		case r == '\'' || r == '"':
			inQuote = !inQuote
		default:
			currentArg += string(r)
		}
	}

	if currentArg != "" {
		args = append(args, currentArg)
	}

	return args
}

func (r *Repl) AddRegisterCmd() {
	cmd := &cobra.Command{
		Use:   "register",
		Short: "register a user",
		Args:  r.RegisterValidation,
		Run:   r.RegisterRunner,
	}

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) RegisterValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	if len(args) != 1 {
		return fmt.Errorf("invalid command")
	}
	// case insensitive
	userName := strings.ToLower(args[0])
	// input validation
	l := len(userName)
	if l < 3 || l > 20 {
		return fmt.Errorf("the [%s] invalid length", userName)
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	if !re.MatchString(userName) {
		return fmt.Errorf("the [%s] contain invalid chars", userName)
	}
	exist := r.storage.IsExistUser(userName)
	if exist {
		return fmt.Errorf("the [%s] has already existed", userName)
	}
	return nil
}

func (r *Repl) RegisterRunner(cmd *cobra.Command, args []string) {
	// case insensitive
	userName := strings.ToLower(args[0])
	r.storage.AddUser(strings.ToLower(args[0]))
	fmt.Printf("Add [%s] successfully\n", userName)
}

func (r *Repl) AddCreateFolderCmd() {
	cmd := &cobra.Command{
		Use:   "create-folder",
		Short: "create a folder for a user",
		Args:  r.CreateFolderValidation,
		Run:   r.CreateFolderRunner,
	}

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) CreateFolderValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)

	switch l {
	case 2, 3:
		// case insensitive
		userName := strings.ToLower(args[0])
		folderName := strings.ToLower(args[1])
		// input validation
		exist := r.storage.IsExistUser(userName)
		if !exist {
			return fmt.Errorf("the [%s] doesn't exist", userName)
		}
		fl := len(folderName)
		if fl < 1 || fl > 100 {
			return fmt.Errorf("the [%s] invalid length", folderName)
		}
		re := regexp.MustCompile(`^[a-zA-Z0-9\.\-\~\_\=\:]+$`)
		if !re.MatchString(folderName) {
			return fmt.Errorf("the [%s] contain invalid chars", folderName)
		}
		exist = r.storage.IsExistFolder(userName, folderName)
		if exist {
			return fmt.Errorf("the [%s] has already existed", folderName)
		}
		if l == 3 && len(args[2]) > 500 {
			return fmt.Errorf("the [description] invalid length")
		}
	default:
		return fmt.Errorf("invalid command")
	}

	return nil
}

func (r *Repl) CreateFolderRunner(cmd *cobra.Command, args []string) {
	// case insensitive
	userName := strings.ToLower(args[0])
	folderName := strings.ToLower(args[1])
	desc := ""
	if len(args) == 3 {
		desc = args[2]
	}

	r.storage.AddFolder(userName, folderName, desc)
	fmt.Printf("Create %s successfully\n", folderName)
}

func (r *Repl) AddDeleteFolderCmd() {
	cmd := &cobra.Command{
		Use:   "delete-folder",
		Short: "delete a folder for a user",
		Args:  r.DeleteFolderValidation,
		Run:   r.DeleteFolderRunner,
	}

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) DeleteFolderValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)
	if l != 2 {
		return fmt.Errorf("invalid command")
	}
	// case insensitive
	userName := strings.ToLower(args[0])
	folderName := strings.ToLower(args[1])
	// input validation
	exist := r.storage.IsExistUser(userName)
	if !exist {
		return fmt.Errorf("the [%s] doesn't exist", userName)
	}
	exist = r.storage.IsExistFolder(userName, folderName)
	if !exist {
		return fmt.Errorf("the [%s] doesn't exist", folderName)

	}

	return nil
}

func (r *Repl) DeleteFolderRunner(cmd *cobra.Command, args []string) {
	// case insensitive
	userName := strings.ToLower(args[0])
	folderName := strings.ToLower(args[1])

	r.storage.DeleteFolder(userName, folderName)
	fmt.Printf("Delete %s successfully\n", folderName)
}

func (r *Repl) AddListFolderCmd() {
	cmd := &cobra.Command{
		Use:   "list-folders",
		Short: "list user folders",
		Args:  r.ListFoldersValidation,
		RunE:  r.ListFoldersRunner,
	}
	cmd.Flags().StringVar(&r.folderSortName, "sort-name", "asc", "Sort by name with asc or desc")
	cmd.Flags().StringVar(&r.folderSortCreated, "sort-created", "", "Sort by created with asc or desc")

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) ListFoldersValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)
	if l != 1 {
		return fmt.Errorf("invalid command")
	}
	// case insensitive
	userName := strings.ToLower(args[0])
	// input validation
	exist := r.storage.IsExistUser(userName)
	if !exist {
		return fmt.Errorf("the [%s] doesn't exist", userName)
	}

	return nil
}

func (r *Repl) ListFoldersRunner(cmd *cobra.Command, args []string) error {
	defer func() {
		r.folderSortName = "asc"
		r.folderSortCreated = ""
	}()

	// case insensitive
	userName := strings.ToLower(args[0])
	sortName := "name"
	orderBy := r.folderSortName
	if r.folderSortCreated != "" {
		sortName = "create"
		orderBy = r.folderSortCreated
	}
	switch orderBy {
	case "asc", "desc":
	default:
		return fmt.Errorf("the [asc|desc] invalid")
	}
	data := r.storage.ListFolder(userName, sortName, strings.ToLower(orderBy))
	for _, v := range data {
		tt := time.Unix(v.FolderCreateTime, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("%s %s %s %s\n", v.FolderName, v.FolderDesc, tt, v.UserName)
	}
	if len(data) == 0 {
		return fmt.Errorf("the [%s] doesn't have any folders", userName)
	}

	return nil
}

func (r *Repl) AddRenameFolderCmd() {
	cmd := &cobra.Command{
		Use:   "rename-folder",
		Short: "rename a folder for a user",
		Args:  r.RenameFolderValidation,
		Run:   r.RenameFolderRunner,
	}

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) RenameFolderValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)

	if l != 3 {
		return fmt.Errorf("invalid command")
	}
	// case insensitive
	userName := strings.ToLower(args[0])
	folderName := strings.ToLower(args[1])
	newFolderName := strings.ToLower(args[2])
	// input validation
	exist := r.storage.IsExistUser(userName)
	if !exist {
		return fmt.Errorf("the [%s] doesn't exist", userName)
	}
	exist = r.storage.IsExistFolder(userName, folderName)
	if !exist {
		return fmt.Errorf("the [%s] doesn't exist", folderName)
	}
	exist = r.storage.IsExistFolder(userName, newFolderName)
	if exist {
		return fmt.Errorf("the [%s] has already existed", newFolderName)
	}
	fl := len(newFolderName)
	if fl < 1 || fl > 100 {
		return fmt.Errorf("the [%s] invalid length", newFolderName)
	}
	re := regexp.MustCompile(`^[a-zA-Z0-9\.\-\~\_\=\:]+$`)
	if !re.MatchString(newFolderName) {
		return fmt.Errorf("the [%s] contain invalid chars", newFolderName)
	}

	return nil
}

func (r *Repl) RenameFolderRunner(cmd *cobra.Command, args []string) {
	// case insensitive
	userName := strings.ToLower(args[0])
	folderName := strings.ToLower(args[1])
	newFolderName := strings.ToLower(args[2])
	r.storage.RenameFolder(userName, folderName, newFolderName)
	fmt.Printf("Rename %s to %s successfully\n", folderName, newFolderName)
}
