package storage

import (
	"fmt"
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

func (t *TestVirtualFileSysStorage) TestDeleteFolders() {
	t.TestStorage.AddUser("test")
	n := 10
	for i := 1; i <= n; i++ {
		t.TestStorage.AddFolder("test", fmt.Sprintf("folder%d", i), "desc")
	}
	for j := n; j > 0; j-- {
		t.TestStorage.DeleteFolder("test", fmt.Sprintf("folder%d", j))
		t.False(t.TestStorage.IsExistFolder("test", fmt.Sprintf("folder%d", j)))
		folders := t.TestStorage.ListFolder("test", "test", "desc")
		t.Equal(j-1, len(folders))
	}
}

func (t *TestVirtualFileSysStorage) TestRenameFolder() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.RenameFolder("test", "folder", "newFolder")
	t.False(t.TestStorage.IsExistFolder("test", "folder"))
	t.True(t.TestStorage.IsExistFolder("test", "newFolder"))
}

func (t *TestVirtualFileSysStorage) TestRenameFolderWithFile() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	t.TestStorage.AddFile("test", "folder", "file", "desc")
	t.TestStorage.RenameFolder("test", "folder", "newFolder")
	t.False(t.TestStorage.IsExistFolder("test", "folder"))
	t.True(t.TestStorage.IsExistFolder("test", "newFolder"))
	t.True(t.TestStorage.IsExistFile("test", "newFolder", "file"))
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

func (t *TestVirtualFileSysStorage) TestDeleteFiles() {
	t.TestStorage.AddUser("test")
	t.TestStorage.AddFolder("test", "folder", "desc")
	n := 10
	for i := 1; i <= n; i++ {
		t.TestStorage.AddFile("test", "folder", fmt.Sprintf("file%d", i), "desc")
	}
	for j := n; j > 0; j-- {
		t.TestStorage.DeleteFile("test", "folder", fmt.Sprintf("file%d", j))
		t.False(t.TestStorage.IsExistFile("test", "folder", fmt.Sprintf("file%d", j)))
		files := t.TestStorage.ListFile("test", "folder", "name", "asc")
		t.Equal(j-1, len(files))
	}
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

func BenchmarkAddUser(b *testing.B) {
	storage := &VirtualFileSysStorage{
		Data:      make(map[string][]VirtualFileSysEntity),
		FolderMap: make(map[string]bool),
		FileMap:   make(map[string]bool),
	}
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("test%d", i)
		storage.AddUser(name)
	}
}

func BenchmarkAddFolder(b *testing.B) {
	storage := &VirtualFileSysStorage{
		Data:      make(map[string][]VirtualFileSysEntity),
		FolderMap: make(map[string]bool),
		FileMap:   make(map[string]bool),
	}
	storage.AddUser("test")
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("folder%d", i)
		desc := "desc01234567890123456789012345678901234567890123456789012345678901234567890123456789"
		storage.AddFolder("test", name, desc)
	}
}

func BenchmarkAddFile(b *testing.B) {
	storage := &VirtualFileSysStorage{
		Data:      make(map[string][]VirtualFileSysEntity),
		FolderMap: make(map[string]bool),
		FileMap:   make(map[string]bool),
	}
	storage.AddUser("test")
	storage.AddFolder("test", "folder", "desc")
	for i := 0; i < b.N; i++ {
		name := fmt.Sprintf("file%d", i)
		desc := "desc01234567890123456789012345678901234567890123456789012345678901234567890123456789"
		storage.AddFile("test", "folder", name, desc)
	}
}

func BenchmarkRenameFolder(b *testing.B) {
	storage := &VirtualFileSysStorage{
		Data:      make(map[string][]VirtualFileSysEntity),
		FolderMap: make(map[string]bool),
		FileMap:   make(map[string]bool),
	}
	storage.AddUser("test")
	folderName := "folder"
	newFolderName := "newfolder"
	storage.AddFolder("test", folderName, "desc")
	for i := 0; i < b.N; i++ {
		if i%2 == 0 {
			storage.RenameFolder("test", folderName, newFolderName)
		} else {
			storage.RenameFolder("test", newFolderName, folderName)
		}
	}
}
