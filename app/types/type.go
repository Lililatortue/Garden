package types

import (
	"time"
)

type Repository struct {
	ID     int64
	Name   string
	UserID int64
	Tags   []*GardenTag
}

type GardenTag struct {
	ID        int64
	Name      string
	Signature string
	Message   string
	Timestamp time.Time
	Tree      *HashTree
}

func NewGardenTag(signature string, message string, timestamp time.Time) *GardenTag {
	return &GardenTag{
		Signature: signature,
		Message:   message,
		Timestamp: timestamp,
	}
}

type FolderNode struct {
	ID        int64
	Signature string
	Filename  string
	Path      string
	Contents  *struct {
		SubFolders []*FolderNode
		SubFiles   []*FileNode
	}
}

type FileNode struct {
	ID        int64
	Signature string
	Filename  string
	Path      string
	Content   string
}

type HashNode interface {
	FileNode | FolderNode
}
type HashTree struct {
	*FolderNode
}
