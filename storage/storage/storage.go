package storage

import (
	"fmt"
)

type InputRes struct {
	Key string
	Url string
}

type Storage interface {
	PutImage(data []byte, contentType string) (*InputRes, error)
	PutVideo(data []byte, contentType string) (*InputRes, error)
	PutObject(key string, data []byte, contentType string) error
	CopyImage(originUrl string) (*InputRes, error)
	GetUrl(key string) (string, error)
	DeleteObject(key string) error
}

var disks = make(map[DiskName]Storage)

func Register(name DiskName, disk Storage) {
	if disk == nil {
		panic("storage: Register disk is nil")
	}
	disks[name] = disk
}

func Disk(name DiskName) (Storage, error) {
	disk, exist := disks[name]
	if !exist {
		return nil, fmt.Errorf("storage: Unknown disk %q", name)
	}
	return disk, nil
}
