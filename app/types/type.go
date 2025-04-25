package types

import (
	"time"
)

type User struct {
	ID           int64         `json:"id" db:"id"`
	Name         string        `json:"name" db:"name"`
	Email        string        `json:"email" db:"email"`
	Password     string        `json:"password" db:"password"`
	Repositories []*Repository `json:"repositories" db:"repositories"`
}

type Repository struct {
	ID       int64     `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	UserID   int64     `json:"user_id" db:"user_id"`
	Branches []*Branch `json:"branches" db:"branches"`
}

type Branch struct {
	ID   int64      `json:"id" db:"id"`
	Name string     `json:"name" db:"name"`
	Head *GardenTag `json:"head" db:"head"`
}

type GardenTag struct {
	ID        int64      `json:"id" db:"id"`
	Parent    *GardenTag `json:"parent" db:"parent"`
	Name      string     `json:"name" db:"name"`
	Signature string     `json:"signature" db:"signature"`
	Message   string     `json:"message" db:"message"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	Tree      *HashTree  `json:"tree" db:"tree"`
}

func NewGardenTag(signature string, message string, timestamp time.Time) *GardenTag {
	return &GardenTag{
		Signature: signature,
		Message:   message,
		Timestamp: timestamp,
	}
}

type FolderNode struct {
	ID        int64  `json:"id" db:"id"`
	Signature string `json:"signature" db:"signature"`
	Filename  string `json:"filename" db:"filename"`
	Path      string `json:"path" db:"path"`
	Contents  *struct {
		SubFolders []*FolderNode `json:"subfolders" db:"subfolders"`
		SubFiles   []*FileNode   `json:"subfiles" db:"subfiles"`
	} `json:"contents" db:"contents"`
}

type FileNode struct {
	ID        int64  `json:"id" db:"id"`
	Signature string `json:"signature" db:"signature"`
	Filename  string `json:"filename" db:"filename"`
	Path      string `json:"path" db:"path"`
	Content   string `json:"content" db:"content"`
}

type HashNode interface {
	FileNode | FolderNode
}
type HashTree struct {
	*FolderNode `json:"folder_node" db:"folder_node"`
}
