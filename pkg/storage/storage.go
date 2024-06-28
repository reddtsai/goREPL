package storage

import (
	"sync"
)

type IStorage interface {
}

type VirtualFileSysStorage struct {
	mu   sync.Mutex
	Data []VirtualFileSysEntity
}

type VirtualFileSysEntity struct {
	UserName         string
	FolderName       string
	FolderCreateTime int64
	FolderDesc       string
	FileName         string
	FileCreateTime   int64
	FileDesc         string
}
