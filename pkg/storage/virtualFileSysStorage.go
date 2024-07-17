//go:generate mockgen -source=storage.go -destination=mock/mock_storage.go -package=mock
package storage

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type VirtualFileSysStorage struct {
	mu        sync.RWMutex
	Data      map[string][]VirtualFileSysEntity
	FolderMap map[string]bool
	FileMap   map[string]bool
}

type VirtualFileSysEntity struct {
	UserName         string
	FolderName       string
	FolderCreateTime int64
	FolderDesc       string
	Files            []VirtualFileSysFileEntity
}

type VirtualFileSysFileEntity struct {
	FileName       string
	FileCreateTime int64
	FileDesc       string
}

func NewVirtualFileSysStorage() IStorage {
	return &VirtualFileSysStorage{
		Data:      make(map[string][]VirtualFileSysEntity),
		FolderMap: make(map[string]bool),
		FileMap:   make(map[string]bool),
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
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].FolderName < entities[j].FolderName
	})
	// 使用二分查找找到第一个符合条件的索引
	index := sort.Search(len(entities), func(i int) bool {
		return entities[i].FolderName >= folderName
	})
	if index < len(entities) && entities[index].FolderName == folderName {
		start := index
		end := index
		end++
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

func (v *VirtualFileSysStorage) RenameFolder(userName, folderName, newFolderName string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	entities := v.Data[userName]
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].FolderName < entities[j].FolderName
	})
	index := sort.Search(len(entities), func(i int) bool {
		return entities[i].FolderName >= folderName
	})
	if index < len(entities) && entities[index].FolderName == folderName {
		entities[index].FolderName = newFolderName
	}
	key := fmt.Sprintf("%s:%s", userName, folderName)
	delete(v.FolderMap, key)
	key = fmt.Sprintf("%s:%s", userName, newFolderName)
	v.FolderMap[key] = true
	for k := range v.FileMap {
		uf := fmt.Sprintf("%s:%s:", userName, folderName)
		if strings.Contains(k, uf) {
			newKey := fmt.Sprintf("%s:%s:%s", userName, newFolderName, k[len(userName)+len(folderName)+2:])
			fmt.Println(newKey, k)
			v.FileMap[newKey] = true
			delete(v.FileMap, k)
		}
	}
}

func (v *VirtualFileSysStorage) IsExistFile(userName, folderName, fileName string) bool {
	v.mu.RLock()
	defer v.mu.RUnlock()

	key := fmt.Sprintf("%s:%s:%s", userName, folderName, fileName)
	_, ok := v.FileMap[key]
	return ok
}

func (v *VirtualFileSysStorage) AddFile(userName, folderName, fileName, fileDesc string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	key := fmt.Sprintf("%s:%s:%s", userName, folderName, fileName)
	v.FileMap[key] = true
	entities := v.Data[userName]
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].FolderName < entities[j].FolderName
	})
	index := sort.Search(len(entities), func(i int) bool {
		return entities[i].FolderName >= folderName
	})
	if index < len(entities) && entities[index].FolderName == folderName {
		entities[index].Files = append(entities[index].Files, VirtualFileSysFileEntity{
			FileName:       fileName,
			FileCreateTime: time.Now().Unix(),
			FileDesc:       fileDesc,
		})
	}
}

func (v *VirtualFileSysStorage) DeleteFile(userName, folderName, fileName string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	entities := v.Data[userName]
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].FolderName < entities[j].FolderName
	})
	index := sort.Search(len(entities), func(i int) bool {
		return entities[i].FolderName >= folderName
	})
	if index < len(entities) && entities[index].FolderName == folderName {
		files := entities[index].Files
		sort.Slice(files, func(i, j int) bool {
			return files[i].FileName < files[j].FileName
		})
		fileIndex := sort.Search(len(files), func(i int) bool {
			return files[i].FileName >= fileName
		})
		if fileIndex < len(files) && files[fileIndex].FileName == fileName {
			start := fileIndex
			end := fileIndex
			end++
			entities[index].Files = append(files[:start], files[end:]...)
		}
	}
	key := fmt.Sprintf("%s:%s:%s", userName, folderName, fileName)
	delete(v.FileMap, key)
}

func (v *VirtualFileSysStorage) ListFile(userName, folderName, sortName, orderBy string) []VirtualFileSysFileEntity {
	v.mu.RLock()
	defer v.mu.RUnlock()

	entities := v.Data[userName]
	sort.Slice(entities, func(i, j int) bool {
		return entities[i].FolderName < entities[j].FolderName
	})
	index := sort.Search(len(entities), func(i int) bool {
		return entities[i].FolderName >= folderName
	})
	if index < len(entities) && entities[index].FolderName == folderName {
		files := entities[index].Files
		switch sortName {
		case "name":
			sort.Slice(files, func(i, j int) bool {
				if orderBy == "desc" {
					return files[i].FileName > files[j].FileName
				}
				return files[i].FileName < files[j].FileName
			})
		case "create":
			sort.Slice(files, func(i, j int) bool {
				if orderBy == "desc" {
					return files[i].FileCreateTime > files[j].FileCreateTime
				}
				return files[i].FileCreateTime < files[j].FileCreateTime
			})
		}
		return files
	}
	return []VirtualFileSysFileEntity{}
}
