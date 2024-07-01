package storage

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TestVirtualFileSysStorage struct {
	suite.Suite

	TestStorage *VirtualFileSysStorage
}

func TestVirtualFileSysStorageSuite(t *testing.T) {
	suite.Run(t, new(TestVirtualFileSysStorage))
}

func (t *TestVirtualFileSysStorage) SetupSuite() {
	t.TestStorage = &VirtualFileSysStorage{
		Data:      make(map[string][]VirtualFileSysEntity),
		FolderMap: make(map[string]bool),
		FileMap:   make(map[string]bool),
	}
}

func (t *TestVirtualFileSysStorage) TestAddUser() {
	t.TestStorage.AddUser("test")
	t.True(t.TestStorage.IsExistUser("test"))
}

func (t *TestVirtualFileSysStorage) TestAddFolder() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.True(t.TestStorage.IsExistFolder("test", "folder"))
}

func (t *TestVirtualFileSysStorage) TestDeleteFolder() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.DeleteFolder("test", "folder")
	t.False(t.TestStorage.IsExistFolder("test", "folder"))
}

func (t *TestVirtualFileSysStorage) TestRenameFolder() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.RenameFolder("test", "folder", "newFolder")
	t.False(t.TestStorage.IsExistFolder("test", "folder"))
	t.True(t.TestStorage.IsExistFolder("test", "newFolder"))
}

func (t *TestVirtualFileSysStorage) TestListFolder() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.AddFolder("test", "folder2", "desc")
	t.TestStorage.AddFolder("test", "folder3", "desc")
	folders := t.TestStorage.ListFolder("test", "name", "asc")
	t.Equal(3, len(folders))
}

func (t *TestVirtualFileSysStorage) TestAddFile() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.AddFile("test", "folder", "file", "desc")
	t.True(t.TestStorage.IsExistFile("test", "folder", "file"))
}

func (t *TestVirtualFileSysStorage) TestDeleteFile() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.AddFile("test", "folder", "file", "desc")
	t.TestStorage.DeleteFile("test", "folder", "file")
	t.False(t.TestStorage.IsExistFile("test", "folder", "file"))
}

func (t *TestVirtualFileSysStorage) TestListFile() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.AddFile("test", "folder", "file", "desc")
	t.TestStorage.AddFile("test", "folder", "file2", "desc")
	t.TestStorage.AddFile("test", "folder", "file3", "desc")
	files := t.TestStorage.ListFile("test", "folder", "name", "asc")
	t.Equal(3, len(files))
}
