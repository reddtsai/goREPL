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

func (t *TestVirtualFileSysStorage) TestListFolderByName() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder1", "desc1")
	t.TestStorage.AddFolder("test", "folder2", "desc2")
	t.TestStorage.AddFolder("test", "folder3", "desc3")
	folders := t.TestStorage.ListFolder("test", "name", "asc")
	t.Equal(3, len(folders))
	t.Equal("folder1", folders[0].FolderName)
}

func (t *TestVirtualFileSysStorage) TestListFolderByNameDesc() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder1", "desc1")
	t.TestStorage.AddFolder("test", "folder2", "desc2")
	t.TestStorage.AddFolder("test", "folder3", "desc3")
	folders := t.TestStorage.ListFolder("test", "name", "desc")
	t.Equal(3, len(folders))
	t.Equal("folder3", folders[0].FolderName)
}

func (t *TestVirtualFileSysStorage) TestListFolderByCreate() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder1", "desc1")
	t.TestStorage.AddFolder("test", "folder2", "desc2")
	t.TestStorage.AddFolder("test", "folder3", "desc3")
	folders := t.TestStorage.ListFolder("test", "create", "asc")
	t.Equal(3, len(folders))
	t.Equal("folder1", folders[0].FolderName)
}

func (t *TestVirtualFileSysStorage) TestListFolderByCreateDesc() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder1", "desc1")
	t.TestStorage.AddFolder("test", "folder2", "desc2")
	t.TestStorage.AddFolder("test", "folder3", "desc3")
	folders := t.TestStorage.ListFolder("test", "create", "desc")
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

func (t *TestVirtualFileSysStorage) TestListFileByNmae() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.AddFile("test", "folder", "file1", "desc1")
	t.TestStorage.AddFile("test", "folder", "file2", "desc2")
	t.TestStorage.AddFile("test", "folder", "file3", "desc3")
	files := t.TestStorage.ListFile("test", "folder", "name", "asc")
	t.Equal(3, len(files))
	t.Equal("file1", files[0].FileName)
}

func (t *TestVirtualFileSysStorage) TestListFileByNmaeDesc() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.AddFile("test", "folder", "file1", "desc1")
	t.TestStorage.AddFile("test", "folder", "file2", "desc2")
	t.TestStorage.AddFile("test", "folder", "file3", "desc3")
	files := t.TestStorage.ListFile("test", "folder", "name", "desc")
	t.Equal(3, len(files))
	t.Equal("file3", files[0].FileName)
}

func (t *TestVirtualFileSysStorage) TestListFileByCreate() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.AddFile("test", "folder", "file1", "desc1")
	t.TestStorage.AddFile("test", "folder", "file2", "desc2")
	t.TestStorage.AddFile("test", "folder", "file3", "desc3")
	files := t.TestStorage.ListFile("test", "folder", "create", "asc")
	t.Equal(3, len(files))
	t.Equal("file1", files[0].FileName)
}

func (t *TestVirtualFileSysStorage) TestListFileByCreateDesc() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.AddFile("test", "folder", "file2", "desc1")
	t.TestStorage.AddFile("test", "folder", "file2", "desc2")
	t.TestStorage.AddFile("test", "folder", "file3", "desc3")
	files := t.TestStorage.ListFile("test", "folder", "create", "desc")
	t.Equal(3, len(files))
}
