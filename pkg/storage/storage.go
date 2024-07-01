package storage

type IStorage interface {
	AddUser(string)
	IsExistUser(string) bool

	AddFolder(string, string, string)
	DeleteFolder(string, string)
	RenameFolder(string, string, string)
	IsExistFolder(string, string) bool
	ListFolder(string, string, string) []VirtualFileSysEntity

	IsExistFile(string, string, string) bool
	AddFile(string, string, string, string)
	DeleteFile(string, string, string)
	ListFile(string, string, string, string) []VirtualFileSysFileEntity
}
