package storage

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type IStorage interface {
	AddUser(string)
	IsExistUser(string) bool

	AddFolder(string, string, string)
	DeleteFolder(string, string)
	IsExistFolder(string, string) bool
	ListFolder(string, string, string) []VirtualFileSysEntity
}

type VirtualFileSysStorage struct {
	mu        sync.RWMutex
	Data      map[string][]VirtualFileSysEntity
	FolderMap map[string]bool
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

type VirtualFileSysEntitiesByFolderName []VirtualFileSysEntity

func (e VirtualFileSysEntitiesByFolderName) Len() int {
	return len(e)
}

func (e VirtualFileSysEntitiesByFolderName) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e VirtualFileSysEntitiesByFolderName) Less(i, j int) bool {
	return e[i].FolderName < e[j].FolderName
}

func NewVirtualFileSysStorage() IStorage {
	return &VirtualFileSysStorage{
		Data:      make(map[string][]VirtualFileSysEntity),
		FolderMap: make(map[string]bool),
	}
}

func (v *VirtualFileSysStorage) AddUser(userName string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.Data[userName] = []VirtualFileSysEntity{}
}

func (v *VirtualFileSysStorage) IsExistUser(userName string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	_, ok := v.Data[userName]
	return ok
}

func (v *VirtualFileSysStorage) AddFolder(userName, folderName, folderDesc string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	key := fmt.Sprintf("%s:%s", userName, folderName)
	v.FolderMap[key] = true
	v.Data[userName] = append(v.Data[userName], VirtualFileSysEntity{
		UserName:         userName,
		FolderName:       folderName,
		FolderCreateTime: time.Now().Unix(),
		FolderDesc:       folderDesc,
	})
}

func (v *VirtualFileSysStorage) IsExistFolder(userName, folderName string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	key := fmt.Sprintf("%s:%s", userName, folderName)
	_, ok := v.FolderMap[key]
	return ok
}

func (v *VirtualFileSysStorage) DeleteFolder(userName, folderName string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	entities := v.Data[userName]
	sort.Sort(VirtualFileSysEntitiesByFolderName(entities))
	// 使用二分查找找到第一个符合条件的索引
	index := sort.Search(len(entities), func(i int) bool {
		return entities[i].FolderName >= folderName
	})
	if index < len(entities) && entities[index].FolderName == folderName {
		start := index
		end := index
		end++
		// 删除目标元素
		v.Data[userName] = append(entities[:start], entities[end:]...)
	}
	key := fmt.Sprintf("%s:%s", userName, folderName)
	delete(v.FolderMap, key)
}

func (v *VirtualFileSysStorage) ListFolder(userName, sortName, orderBy string) []VirtualFileSysEntity {
	v.mu.RLock()
	defer v.mu.RUnlock()
	entities := v.Data[userName]

	switch sortName {
	case "name":
		sort.Slice(entities, func(i, j int) bool {
			if orderBy == "desc" {
				return entities[i].FolderName > entities[j].FolderName
			}
			return entities[i].FolderName < entities[j].FolderName
		})
	case "create":
		sort.Slice(entities, func(i, j int) bool {
			if orderBy == "desc" {
				return entities[i].FolderCreateTime > entities[j].FolderCreateTime
			}
			return entities[i].FolderCreateTime < entities[j].FolderCreateTime
		})
	}

	return entities
}
