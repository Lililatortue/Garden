package types

import (
	"time"
)

type User struct {
	ID           int64         `json:"id,omitempty" db:"id"`
	Name         string        `json:"name" db:"name"`
	Email        string        `json:"email" db:"email"`
	Password     string        `json:"password" db:"password"`
	Repositories []*Repository `json:"repositories,omitempty" db:"repositories"`
}

func NewUser(opts ...func(*User)) *User {
	u := DefaultUser
	for _, opt := range opts {
		opt(&u)
	}

	return &u
}

type Repository struct {
	ID       int64     `json:"id,omitempty" db:"id"`
	Name     string    `json:"name" db:"name"`
	UserID   int64     `json:"user_id" db:"user_id"`
	Branches []*Branch `json:"branches,omitempty" db:"branches"`
}

type Branch struct {
	ID   int64     `json:"id,omitempty" db:"id"`
	Name string    `json:"name" db:"name"`
	Head GardenTag `json:"head,omitempty" db:"head"`
}

func NewBranch(opts ...func(*Branch)) *Branch {
	branch := DefaultBranch

	for _, opt := range opts {
		opt(&branch)
	}

	return &branch
}

type GardenTag struct {
	ID        int64      `json:"id,omitempty" db:"id"`
	Parent    *GardenTag `json:"parent,omitempty" db:"parent"`
	Name      string     `json:"name" db:"name"`
	Signature string     `json:"signature" db:"signature"`
	Message   string     `json:"message" db:"message"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	Tree      HashTree   `json:"tree" db:"tree"`
}

type List[T any] []*T

type FolderNode struct {
	ID         int64            `json:"id,omitempty" db:"id"`
	Signature  string           `json:"signature" db:"signature"`
	Filename   string           `json:"filename" db:"filename"`
	Path       string           `json:"path,omitempty" db:"path"`
	SubFolders List[FolderNode] `json:"subfolders,omitempty" db:"subfolders"`
	SubFiles   List[FileNode]   `json:"subfiles,omitempty" db:"subfiles"`
}

type FileNode struct {
	ID        int64  `json:"id,omitempty" db:"id"`
	Signature string `json:"signature" db:"signature"`
	Filename  string `json:"filename" db:"filename"`
	Path      string `json:"path,omitempty" db:"path"`
	Content   string `json:"content" db:"content"`
}

func NewFileNode(opts ...func(*FileNode)) *FileNode {
	file := DefaultFileNode

	for _, opt := range opts {
		opt(&file)
	}

	return &file
}

type HashNode interface {
	FileNode | FolderNode
}
type HashTree struct {
	FolderNode `json:"folder_node" db:"folder_node"`
}
