package storage

type IStorage interface {
	AddUser(userName string)
	IsExistUser(userName string) bool

	AddFolder(userName, folderName, folderDesc string)
	DeleteFolder(userName, folderName string)
	RenameFolder(userName, folderName, newFolderName string)
	IsExistFolder(userName, folderName string) bool
	ListFolder(userName, sortName, orderBy string) []VirtualFileSysEntity

	IsExistFile(userName, folderName, fileName string) bool
	AddFile(userName, folderName, fileName, fileDesc string)
	DeleteFile(userName, folderName, fileName string)
	ListFile(userName, folderName, sortName, orderBy string) []VirtualFileSysFileEntity
}
