package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/reddtsai/goREPL/pkg/storage"
	"github.com/reddtsai/goREPL/pkg/storage/mock"
)

type TestRepl struct {
	suite.Suite
	ctrl *gomock.Controller

	repl        *Repl
	mockStorage *mock.MockIStorage
}

func TestReplSuite(t *testing.T) {
	suite.Run(t, new(TestRepl))
}

func (t *TestRepl) SetupSuite() {
	t.ctrl = gomock.NewController(t.T())
	t.mockStorage = mock.NewMockIStorage(t.ctrl)
	t.repl = &Repl{
		storage: t.mockStorage,
	}
	rootCmd := &cobra.Command{
		Use:     "repl",
		Version: "test",
		Run:     t.repl.RootCmdRunner,
	}
	t.repl.rootCmd = rootCmd
	t.repl.AddRegisterCmd()
	t.repl.AddCreateFolderCmd()
	t.repl.AddDeleteFolderCmd()
	t.repl.AddListFoldersCmd()
	t.repl.AddRenameFolderCmd()
	t.repl.AddCreateFileCmd()
	t.repl.AddDeleteFileCmd()
	t.repl.AddListFilesCmd()
	t.repl.Execute()
}

func (t *TestRepl) TearDownSuite() {
	t.ctrl.Finish()
}

func (t *TestRepl) Execute(args []string) (string, error) {
	sout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	t.repl.rootCmd.SetArgs(args)
	err := t.repl.rootCmd.Execute()
	w.Close()
	os.Stdout = sout
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()

	return buf.String(), err
}

func (t *TestRepl) TestRegisterCmdSuccess() {
	userName := "test"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(false)
	t.mockStorage.EXPECT().AddUser(userName)
	// execute
	out, err := t.Execute([]string{"register", userName})
	// testing
	assert.Nil(t.T(), err)
	expected := fmt.Sprintf("Add [%s] successfully\n", userName)
	assert.Equal(t.T(), expected, out)
}

func (t *TestRepl) TestRegisterCmdUnrecognizedArgs() {
	// execute
	_, err := t.Execute([]string{"register"})
	// testing
	assert.NotNil(t.T(), err)
	cmd, _, _ := t.repl.rootCmd.Find([]string{"register"})
	expected := fmt.Sprintf("unrecognized argument\n%s", cmd.UsageString())
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestRegisterCmdUserNameExist() {
	userName := "test"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	// execute
	_, err := t.Execute([]string{"register", userName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] has already existed", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestRegisterCmdUserNameInvalidChar() {
	userName := "test@123"
	// execute
	_, err := t.Execute([]string{"register", userName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] contain invalid chars", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestRegisterCmdUserNameInvalidLength() {
	userName := "t"
	// execute
	_, err := t.Execute([]string{"register", userName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] invalid length", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFolderCmdSuccess() {
	userName := "test"
	folderName := "folder"
	folderDesc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(false)
	t.mockStorage.EXPECT().AddFolder(userName, folderName, folderDesc)
	// execute
	out, err := t.Execute([]string{"create-folder", userName, folderName, folderDesc})
	// testing
	assert.Nil(t.T(), err)
	expected := fmt.Sprintf("Create [%s] successfully\n", folderName)
	assert.Equal(t.T(), expected, out)
}

func (t *TestRepl) TestCreateFolderCmdUnrecognizedArgs() {
	// execute
	_, err := t.Execute([]string{"create-folder"})
	// testing
	assert.NotNil(t.T(), err)
	cmd, _, _ := t.repl.rootCmd.Find([]string{"create-folder"})
	expected := fmt.Sprintf("unrecognized argument\n%s", cmd.UsageString())
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFolderCmdUserNameNotExist() {
	userName := "test"
	folderName := "folder"
	folderDesc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(false)
	// execute
	_, err := t.Execute([]string{"create-folder", userName, folderName, folderDesc})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFolderCmdFolderNameInvalidLength() {
	userName := "test"
	folderName := "f12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	folderDesc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	// execute
	_, err := t.Execute([]string{"create-folder", userName, folderName, folderDesc})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] invalid length", folderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFolderCmdFolderNameInvalidChar() {
	userName := "test"
	folderName := "folder@123"
	folderDesc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	// execute
	_, err := t.Execute([]string{"create-folder", userName, folderName, folderDesc})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] contain invalid chars", folderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFolderCmdFolderNameExist() {
	userName := "test"
	folderName := "folder"
	folderDesc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	// execute
	_, err := t.Execute([]string{"create-folder", userName, folderName, folderDesc})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] has already existed", folderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFolderCmdFolderDescInvalidLength() {
	userName := "test"
	folderName := "folder"
	folderDesc := "desc1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(false)
	// execute
	_, err := t.Execute([]string{"create-folder", userName, folderName, folderDesc})
	// testing
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), "the [description] invalid length", err.Error())
}

func (t *TestRepl) TestDeleteFolderCmdSuccess() {
	userName := "test"
	folderName := "folder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().DeleteFolder(userName, folderName)
	// execute
	out, err := t.Execute([]string{"delete-folder", userName, folderName})
	// testing
	assert.Nil(t.T(), err)
	expected := fmt.Sprintf("Delete [%s] successfully\n", folderName)
	assert.Equal(t.T(), expected, out)
}

func (t *TestRepl) TestDeleteFolderCmdUnrecognizedArgs() {
	// execute
	_, err := t.Execute([]string{"delete-folder"})
	// testing
	assert.NotNil(t.T(), err)
	cmd, _, _ := t.repl.rootCmd.Find([]string{"delete-folder"})
	expected := fmt.Sprintf("unrecognized argument\n%s", cmd.UsageString())
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestDeleteFolderCmdUserNameNotExist() {
	userName := "test"
	folderName := "folder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(false)
	// execute
	_, err := t.Execute([]string{"delete-folder", userName, folderName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestDeleteFolderCmdFolderNameNotExist() {
	userName := "test"
	folderName := "folder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(false)
	// execute
	_, err := t.Execute([]string{"delete-folder", userName, folderName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", folderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestListFoldersCmdSuccess() {
	userName := "test"
	folders := []storage.VirtualFileSysEntity{
		{UserName: "test", FolderName: "folder1", FolderCreateTime: 1719797050, FolderDesc: "desc1"},
		{UserName: "test", FolderName: "folder2", FolderCreateTime: 719797050, FolderDesc: "desc2"},
		{UserName: "test", FolderName: "folder3", FolderCreateTime: 1719797050, FolderDesc: "desc3"},
	}
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().ListFolder(userName, "name", "asc").Return(folders)
	// execute
	out, err := t.Execute([]string{"list-folders", userName})
	// testing
	assert.Nil(t.T(), err)
	list := strings.Split(out, "\n")
	assert.Equal(t.T(), 3+1, len(list))
	assert.Contains(t.T(), list[2], "folder3")
}

func (t *TestRepl) TestListFoldersCmdByNameSuccess() {
	userName := "test"
	folders := []storage.VirtualFileSysEntity{
		{UserName: "test", FolderName: "folder3", FolderCreateTime: 1719797050, FolderDesc: "desc3"},
		{UserName: "test", FolderName: "folder2", FolderCreateTime: 719797050, FolderDesc: "desc2"},
		{UserName: "test", FolderName: "folder1", FolderCreateTime: 1719797050, FolderDesc: "desc1"},
	}
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().ListFolder(userName, "name", "desc").Return(folders)
	// execute
	out, err := t.Execute([]string{"list-folders", userName, "--sort-name", "desc"})
	// testing
	assert.Nil(t.T(), err)
	list := strings.Split(out, "\n")
	assert.Equal(t.T(), 3+1, len(list))
	assert.Contains(t.T(), list[0], "folder3")
}

func (t *TestRepl) TestListFoldersCmdByCreateSuccess() {
	userName := "test"
	folders := []storage.VirtualFileSysEntity{
		{UserName: "test", FolderName: "folder3", FolderCreateTime: 1719797053, FolderDesc: "desc3"},
		{UserName: "test", FolderName: "folder2", FolderCreateTime: 719797052, FolderDesc: "desc2"},
		{UserName: "test", FolderName: "folder1", FolderCreateTime: 1719797051, FolderDesc: "desc1"},
	}
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().ListFolder(userName, "create", "desc").Return(folders)
	// execute
	out, err := t.Execute([]string{"list-folders", userName, "--sort-created", "desc"})
	// testing
	assert.Nil(t.T(), err)
	list := strings.Split(out, "\n")
	assert.Equal(t.T(), 3+1, len(list))
	assert.Contains(t.T(), list[0], "folder3")
}

func (t *TestRepl) TestListFoldersCmdNoData() {
	userName := "test"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().ListFolder(userName, "name", "asc").Return([]storage.VirtualFileSysEntity{})
	// execute
	out, err := t.Execute([]string{"list-folders", userName})
	// testing
	assert.Nil(t.T(), err)
	expected := fmt.Sprintf("Warning: the [%s] doesn't have any folders\n", userName)
	assert.Equal(t.T(), expected, out)
}

func (t *TestRepl) TestListFoldersCmdUnrecognizedArgs() {
	// execute
	_, err := t.Execute([]string{"list-folders"})
	// testing
	assert.NotNil(t.T(), err)
	cmd, _, _ := t.repl.rootCmd.Find([]string{"list-folders"})
	expected := fmt.Sprintf("unrecognized argument\n%s", cmd.UsageString())
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestListFoldersCmdUserNameNotExist() {
	userName := "test"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(false)
	// execute
	_, err := t.Execute([]string{"list-folders", userName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestRenameFolderCmdSuccess() {
	userName := "test"
	folderName := "folder"
	newFolderName := "newfolder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true).Times(1)
	t.mockStorage.EXPECT().IsExistFolder(userName, newFolderName).Return(false).Times(1)
	t.mockStorage.EXPECT().RenameFolder(userName, folderName, newFolderName)
	// execute
	out, err := t.Execute([]string{"rename-folder", userName, folderName, newFolderName})
	// testing
	assert.Nil(t.T(), err)
	expected := fmt.Sprintf("Rename [%s] to [%s] successfully\n", folderName, newFolderName)
	assert.Equal(t.T(), expected, out)
}

func (t *TestRepl) TestRenameFolderCmdUnrecognizedArgs() {
	// execute
	_, err := t.Execute([]string{"rename-folder"})
	// testing
	assert.NotNil(t.T(), err)
	cmd, _, _ := t.repl.rootCmd.Find([]string{"rename-folder"})
	expected := fmt.Sprintf("unrecognized argument\n%s", cmd.UsageString())
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestRenameFolderCmdUserNameNotExist() {
	userName := "test"
	folderName := "folder"
	newFolderName := "newFolder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(false)
	// execute
	_, err := t.Execute([]string{"rename-folder", userName, folderName, newFolderName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestRenameFolderCmdFolderNameNotExist() {
	userName := "test"
	folderName := "folder"
	newFolderName := "newfolder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(false)
	// execute
	_, err := t.Execute([]string{"rename-folder", userName, folderName, newFolderName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", folderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestRenameFolderCmdNewFolderNameExist() {
	userName := "test"
	folderName := "folder"
	newFolderName := "newfolder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, newFolderName).Return(true)
	// execute
	_, err := t.Execute([]string{"rename-folder", userName, folderName, newFolderName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] has already existed", newFolderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestRenameFolderCmdNewFolderNameInvalidLength() {
	userName := "test"
	folderName := "folder"
	newFolderName := "newfolder12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, newFolderName).Return(false)
	// execute
	_, err := t.Execute([]string{"rename-folder", userName, folderName, newFolderName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] invalid length", newFolderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestRenameFolderCmdNewFolderNameInvalidChar() {
	userName := "test"
	folderName := "folder"
	newFolderName := "newfolder@123"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, newFolderName).Return(false)
	// execute
	_, err := t.Execute([]string{"rename-folder", userName, folderName, newFolderName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] contain invalid chars", newFolderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFileCmdSuccess() {
	userName := "test"
	folderName := "folder"
	fileName := "file"
	desc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().IsExistFile(userName, folderName, fileName).Return(false)
	t.mockStorage.EXPECT().AddFile(userName, folderName, fileName, desc)
	// execute
	out, err := t.Execute([]string{"create-file", userName, folderName, fileName, desc})
	// testing
	assert.Nil(t.T(), err)
	expected := fmt.Sprintf("Create [%s] in [%s]/[%s] successfully\n", fileName, userName, folderName)
	assert.Equal(t.T(), expected, out)
}

func (t *TestRepl) TestCreateFileCmdUnrecognizedArgs() {
	// execute
	_, err := t.Execute([]string{"create-file"})
	// testing
	assert.NotNil(t.T(), err)
	cmd, _, _ := t.repl.rootCmd.Find([]string{"create-file"})
	expected := fmt.Sprintf("unrecognized argument\n%s", cmd.UsageString())
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFileCmdUserNameNotExist() {
	userName := "test"
	folderName := "folder"
	fileName := "file"
	desc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(false)
	// execute
	_, err := t.Execute([]string{"create-file", userName, folderName, fileName, desc})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFileCmdFolderNameNotExist() {
	userName := "test"
	folderName := "folder"
	fileName := "file"
	desc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(false)
	// execute
	_, err := t.Execute([]string{"create-file", userName, folderName, fileName, desc})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", folderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFileCmdFileNameInvalidLength() {
	userName := "test"
	folderName := "folder"
	fileName := "f12345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	// execute
	_, err := t.Execute([]string{"create-file", userName, folderName, fileName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] invalid length", fileName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFileCmdFileNameInvalidChar() {
	userName := "test"
	folderName := "folder"
	fileName := "file@123"
	desc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	// execute
	_, err := t.Execute([]string{"create-file", userName, folderName, fileName, desc})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] contain invalid chars", fileName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFileCmdFileNameExist() {
	userName := "test"
	folderName := "folder"
	fileName := "file"
	desc := "desc"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().IsExistFile(userName, folderName, fileName).Return(true)
	// execute
	_, err := t.Execute([]string{"create-file", userName, folderName, fileName, desc})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] has already existed", fileName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestCreateFileCmdFileDescInvalidLength() {
	userName := "test"
	folderName := "folder"
	fileName := "file"
	desc := "desc1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().IsExistFile(userName, folderName, fileName).Return(false)
	// execute
	_, err := t.Execute([]string{"create-file", userName, folderName, fileName, desc})
	// testing
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), "the [description] invalid length", err.Error())
}

func (t *TestRepl) TestDeleteFileCmdSuccess() {
	userName := "test"
	folderName := "folder"
	fileName := "file"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().IsExistFile(userName, folderName, fileName).Return(true)
	t.mockStorage.EXPECT().DeleteFile(userName, folderName, fileName)
	// execute
	out, err := t.Execute([]string{"delete-file", userName, folderName, fileName})
	// testing
	assert.Nil(t.T(), err)
	expected := fmt.Sprintf("Delete [%s] in [%s]/[%s] successfully\n", fileName, userName, folderName)
	assert.Equal(t.T(), expected, out)
}

func (t *TestRepl) TestDeleteFileCmdUnrecognizedArgs() {
	// execute
	_, err := t.Execute([]string{"delete-file"})
	// testing
	assert.NotNil(t.T(), err)
	cmd, _, _ := t.repl.rootCmd.Find([]string{"delete-file"})
	expected := fmt.Sprintf("unrecognized argument\n%s", cmd.UsageString())
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestDeleteFileCmdUserNameNotExist() {
	userName := "test"
	folderName := "folder"
	fileName := "file"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(false)
	// execute
	_, err := t.Execute([]string{"delete-file", userName, folderName, fileName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestDeleteFileCmdFolderNameNotExist() {
	userName := "test"
	folderName := "folder"
	fileName := "file"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(false)
	// execute
	_, err := t.Execute([]string{"delete-file", userName, folderName, fileName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", folderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestDeleteFileCmdFileNameNotExist() {
	userName := "test"
	folderName := "folder"
	fileName := "file"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().IsExistFile(userName, folderName, fileName).Return(false)
	// execute
	_, err := t.Execute([]string{"delete-file", userName, folderName, fileName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", fileName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestListFilesCmdSuccess() {
	userName := "test"
	folderName := "folder"
	files := []storage.VirtualFileSysFileEntity{
		{FileName: "file1", FileCreateTime: 1719797050, FileDesc: "desc1"},
		{FileName: "file2", FileCreateTime: 1719797050, FileDesc: "desc2"},
		{FileName: "file3", FileCreateTime: 1719797050, FileDesc: "desc3"},
	}
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().ListFile(userName, folderName, "name", "asc").Return(files)
	// execute
	out, err := t.Execute([]string{"list-files", userName, folderName})
	// testing
	assert.Nil(t.T(), err)
	list := strings.Split(out, "\n")
	assert.Equal(t.T(), 3+1, len(list))
	assert.Contains(t.T(), list[0], "file1")
}

func (t *TestRepl) TestListFilesCmdByNameSuccess() {
	userName := "test"
	folderName := "folder"
	files := []storage.VirtualFileSysFileEntity{
		{FileName: "file3", FileCreateTime: 1719797050, FileDesc: "desc3"},
		{FileName: "file2", FileCreateTime: 1719797050, FileDesc: "desc2"},
		{FileName: "file1", FileCreateTime: 1719797050, FileDesc: "desc1"},
	}
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().ListFile(userName, folderName, "name", "desc").Return(files)
	// execute
	out, err := t.Execute([]string{"list-files", userName, folderName, "--sort-name", "desc"})
	// testing
	assert.Nil(t.T(), err)
	list := strings.Split(out, "\n")
	assert.Equal(t.T(), 3+1, len(list))
	assert.Contains(t.T(), list[0], "file3")
}

func (t *TestRepl) TestListFilesCmdByCreateSuccess() {
	userName := "test"
	folderName := "folder"
	files := []storage.VirtualFileSysFileEntity{
		{FileName: "file3", FileCreateTime: 1719797053, FileDesc: "desc3"},
		{FileName: "file2", FileCreateTime: 1719797052, FileDesc: "desc2"},
		{FileName: "file1", FileCreateTime: 1719797051, FileDesc: "desc1"},
	}
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().ListFile(userName, folderName, "create", "desc").Return(files)
	// execute
	out, err := t.Execute([]string{"list-files", userName, folderName, "--sort-created", "desc"})
	// testing
	assert.Nil(t.T(), err)
	list := strings.Split(out, "\n")
	assert.Equal(t.T(), 3+1, len(list))
	assert.Contains(t.T(), list[0], "file3")
}

func (t *TestRepl) TestListFilesCmdNoData() {
	userName := "test"
	folderName := "folder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(true)
	t.mockStorage.EXPECT().ListFile(userName, folderName, "name", "asc").Return([]storage.VirtualFileSysFileEntity{})
	// execute
	out, err := t.Execute([]string{"list-files", userName, folderName})
	// testing
	assert.Nil(t.T(), err)
	expected := fmt.Sprintf("Warning: the [%s] is empty\n", folderName)
	assert.Equal(t.T(), expected, out)
}

func (t *TestRepl) TestListFilesCmdUnrecognizedArgs() {
	// execute
	_, err := t.Execute([]string{"list-files"})
	// testing
	assert.NotNil(t.T(), err)
	cmd, _, _ := t.repl.rootCmd.Find([]string{"list-files"})
	expected := fmt.Sprintf("unrecognized argument\n%s", cmd.UsageString())
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestListFilesCmdUserNameNotExist() {
	userName := "test"
	folderName := "folder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(false)
	// execute
	_, err := t.Execute([]string{"list-files", userName, folderName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", userName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestListFilesCmdFolderNameNotExist() {
	userName := "test"
	folderName := "folder"
	// mock data
	t.mockStorage.EXPECT().IsExistUser(userName).Return(true)
	t.mockStorage.EXPECT().IsExistFolder(userName, folderName).Return(false)
	// execute
	_, err := t.Execute([]string{"list-files", userName, folderName})
	// testing
	assert.NotNil(t.T(), err)
	expected := fmt.Sprintf("the [%s] doesn't exist", folderName)
	assert.Equal(t.T(), expected, err.Error())
}

func (t *TestRepl) TestHelpCmd() {
	// execute
	t.repl.HelpCmd()
}

func (t *TestRepl) TestSplitArgs() {
	str := "cmd 'new folder' 'hello world'"
	s := t.repl.SplitArgs(str)
	assert.Equal(t.T(), 3, len(s))
}
