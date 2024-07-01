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
	fileSortName      string
	fileSortCreated   string
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
			r.Help()
		default:
			args := r.SplitArgs(line)
			if len(args) > 0 {
				foundCmd, _, err := cmd.Find(args)
				if err != nil || foundCmd == r.rootCmd {
					fmt.Fprintln(os.Stderr, "Error: unrecognized command.")
					continue
				}
				err = foundCmd.ParseFlags(args)
				if err != nil {
					fmt.Println(foundCmd.UsageString())
					continue
				}

				cmd.SetArgs(args)
				_ = cmd.Execute()
			}
		}
	}
}

func (r *Repl) Help() {
	fmt.Println("Usage:")
	fmt.Println("  register [username]")
	fmt.Println("  create-folder [username] [foldername] [description]?")
	fmt.Println("  delete-folder [username] [foldername]")
	fmt.Println("  list-folders [username] [--sort-name|--sort-created] [asc|desc]")
	fmt.Println("  rename-folder [username] [foldername] [new-foldername]")
	fmt.Println("  create-file [username] [foldername] [filename] [description]?")
	fmt.Println("  delete-file [username] [foldername] [filename]")
	fmt.Println("  list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")
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
	cmd.SetUsageTemplate("Usage:\n  register [username]")

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) RegisterValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	if len(args) != 1 {
		return fmt.Errorf("unrecognized argument\n%s", cmd.UsageString())
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
	cmd.SetUsageTemplate("Usage:\n  create-folder [username] [foldername] [description]?")

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
		return fmt.Errorf("unrecognized argument\n%s", cmd.UsageString())
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
	fmt.Printf("Create [%s] successfully\n", folderName)
}

func (r *Repl) AddDeleteFolderCmd() {
	cmd := &cobra.Command{
		Use:   "delete-folder",
		Short: "delete a folder for a user",
		Args:  r.DeleteFolderValidation,
		Run:   r.DeleteFolderRunner,
	}
	cmd.SetUsageTemplate("Usage:\n  delete-folder [username] [foldername]")

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) DeleteFolderValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)
	if l != 2 {
		return fmt.Errorf("unrecognized argument\n%s", cmd.UsageString())
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
	fmt.Printf("Delete [%s] successfully\n", folderName)
}

func (r *Repl) AddListFoldersCmd() {
	cmd := &cobra.Command{
		Use:   "list-folders",
		Short: "list user folders",
		Args:  r.ListFoldersValidation,
		Run:   r.ListFoldersRunner,
	}
	cmd.Flags().StringVar(&r.folderSortName, "sort-name", "", "Sort by name with asc or desc")
	cmd.Flags().StringVar(&r.folderSortCreated, "sort-created", "", "Sort by created with asc or desc")
	cmd.SetUsageTemplate("Usage:\n  list-folders [username] [--sort-name|--sort-created] [asc|desc]")

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) ListFoldersValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)
	if l != 1 {
		return fmt.Errorf("unrecognized argument\n%s", cmd.UsageString())
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

func (r *Repl) ListFoldersRunner(cmd *cobra.Command, args []string) {
	defer func() {
		r.folderSortName = ""
		r.folderSortCreated = ""
	}()

	// case insensitive
	userName := strings.ToLower(args[0])
	sortName, orderBy := "", ""
	if r.folderSortCreated != "" {
		sortName = "create"
		orderBy = r.folderSortCreated
	} else if r.folderSortName != "" {
		sortName = "name"
		orderBy = r.folderSortName
	} else {
		// if neither sort-created nor sort-name
		sortName = "name"
		orderBy = "asc"
	}
	switch orderBy {
	case "asc", "desc":
	default:
		fmt.Println(cmd.UsageString())
		return
	}
	data := r.storage.ListFolder(userName, sortName, strings.ToLower(orderBy))
	for _, v := range data {
		tt := time.Unix(v.FolderCreateTime, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("%s %s %s %s\n", v.FolderName, v.FolderDesc, tt, v.UserName)
	}
	if len(data) == 0 {
		fmt.Printf("Warning: the [%s] doesn't have any folders\n", userName)
	}
}

func (r *Repl) AddRenameFolderCmd() {
	cmd := &cobra.Command{
		Use:   "rename-folder",
		Short: "rename a folder for a user",
		Args:  r.RenameFolderValidation,
		Run:   r.RenameFolderRunner,
	}
	cmd.SetUsageTemplate("Usage:\n  rename-folder [username] [foldername] [new-foldername]")

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) RenameFolderValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)

	if l != 3 {
		return fmt.Errorf("unrecognized argument\n%s", cmd.UsageString())
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
	fmt.Printf("Rename [%s] to [%s] successfully\n", folderName, newFolderName)
}

func (r *Repl) AddCreateFileCmd() {
	cmd := &cobra.Command{
		Use:   "create-file",
		Short: "create a file for a user",
		Args:  r.CreateFileValidation,
		Run:   r.CreateFileRunner,
	}
	cmd.SetUsageTemplate("Usage:\n  create-file [username] [foldername] [filename] [description]?")

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) CreateFileValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)

	switch l {
	case 3, 4:
		// case insensitive
		userName := strings.ToLower(args[0])
		folderName := strings.ToLower(args[1])
		fileName := strings.ToLower(args[2])
		// input validation
		exist := r.storage.IsExistUser(userName)
		if !exist {
			return fmt.Errorf("the [%s] doesn't exist", userName)
		}
		exist = r.storage.IsExistFolder(userName, folderName)
		if !exist {
			return fmt.Errorf("the [%s] doesn't exist", folderName)
		}
		fl := len(fileName)
		if fl < 1 || fl > 100 {
			return fmt.Errorf("the [%s] invalid length", fileName)
		}
		re := regexp.MustCompile(`^[a-zA-Z0-9\.\-\~\_\=\:]+$`)
		if !re.MatchString(fileName) {
			return fmt.Errorf("the [%s] contain invalid chars", fileName)
		}
		exist = r.storage.IsExistFile(userName, folderName, fileName)
		if exist {
			return fmt.Errorf("the [%s] has already existed", fileName)
		}
		if l == 4 && len(args[3]) > 500 {
			return fmt.Errorf("the [description] invalid length")
		}
	default:
		return fmt.Errorf("unrecognized argumen\n%s", cmd.UsageString())
	}

	return nil
}

func (r *Repl) CreateFileRunner(cmd *cobra.Command, args []string) {
	// case insensitive
	userName := strings.ToLower(args[0])
	folderName := strings.ToLower(args[1])
	fileName := strings.ToLower(args[2])
	desc := ""
	if len(args) == 4 {
		desc = args[3]
	}

	r.storage.AddFile(userName, folderName, fileName, desc)
	fmt.Printf("Create [%s] in [%s]/[%s] successfully\n", fileName, userName, folderName)
}

func (r *Repl) AddDeleteFileCmd() {
	cmd := &cobra.Command{
		Use:   "delete-file",
		Short: "delete a file for a user",
		Args:  r.DeleteFileValidation,
		Run:   r.DeleteFileRunner,
	}
	cmd.SetUsageTemplate("Usage:\n  delete-file [username] [foldername] [filename]")

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) DeleteFileValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)
	if l != 3 {
		return fmt.Errorf("unrecognized argument\n%s", cmd.UsageString())
	}
	// case insensitive
	userName := strings.ToLower(args[0])
	folderName := strings.ToLower(args[1])
	fileName := strings.ToLower(args[2])
	// input validation
	exist := r.storage.IsExistUser(userName)
	if !exist {
		return fmt.Errorf("the [%s] doesn't exist", userName)
	}
	exist = r.storage.IsExistFolder(userName, folderName)
	if !exist {
		return fmt.Errorf("the [%s] doesn't exist", folderName)

	}
	exist = r.storage.IsExistFile(userName, folderName, fileName)
	if !exist {
		return fmt.Errorf("the [%s] doesn't exist", fileName)
	}

	return nil
}

func (r *Repl) DeleteFileRunner(cmd *cobra.Command, args []string) {
	// case insensitive
	userName := strings.ToLower(args[0])
	folderName := strings.ToLower(args[1])
	fileName := strings.ToLower(args[2])

	r.storage.DeleteFile(userName, folderName, fileName)
	fmt.Printf("Delete [%s] in [%s]/[%s] successfully\n", fileName, userName, folderName)
}

func (r *Repl) AddListFilesCmd() {
	cmd := &cobra.Command{
		Use:   "list-files",
		Short: "list user files",
		Args:  r.ListFilesValidation,
		Run:   r.ListFilesRunner,
	}
	cmd.Flags().StringVar(&r.fileSortName, "sort-name", "asc", "Sort by name with asc or desc")
	cmd.Flags().StringVar(&r.fileSortCreated, "sort-created", "", "Sort by created with asc or desc")
	cmd.SetUsageTemplate("Usage:\n  list-files [username] [foldername] [--sort-name|--sort-created] [asc|desc]")

	r.rootCmd.AddCommand(cmd)
}

func (r *Repl) ListFilesValidation(cmd *cobra.Command, args []string) error {
	cmd.SilenceUsage = true
	l := len(args)
	if l != 2 {
		return fmt.Errorf("unrecognized argument\n%s", cmd.UsageString())
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

func (r *Repl) ListFilesRunner(cmd *cobra.Command, args []string) {
	defer func() {
		r.fileSortName = ""
		r.fileSortCreated = ""
	}()

	// case insensitive
	userName := strings.ToLower(args[0])
	folderName := strings.ToLower(args[1])
	sortName, orderBy := "", ""
	if r.fileSortCreated != "" {
		sortName = "create"
		orderBy = r.fileSortCreated
	} else if r.folderSortName != "" {
		sortName = "name"
		orderBy = r.folderSortName
	} else {
		// if neither sort-created nor sort-name
		sortName = "name"
		orderBy = "asc"
	}
	switch orderBy {
	case "asc", "desc":
	default:
		fmt.Println(cmd.UsageString())
		return
	}
	data := r.storage.ListFile(userName, folderName, sortName, strings.ToLower(orderBy))
	for _, v := range data {
		tt := time.Unix(v.FileCreateTime, 0).Format("2006-01-02 15:04:05")
		fmt.Printf("%s %s %s %s %s\n", v.FileName, v.FileDesc, tt, folderName, userName)
	}
	if len(data) == 0 {
		fmt.Printf("Warning: the [%s] is empty\n", folderName)
	}
}
