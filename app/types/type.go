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

func NewUser(opts ...func(*User)) User {
	u := User{
		ID:           -1,
		Name:         "",
		Email:        "",
		Password:     "",
		Repositories: []*Repository{},
	}
	for _, opt := range opts {
		opt(&u)
	}

	return u

}

type Repository struct {
	ID       int64     `json:"id" db:"id"`
	Name     string    `json:"name" db:"name"`
	UserID   int64     `json:"user_id" db:"user_id"`
	Branches []*Branch `json:"branches" db:"branches"`
}

func NewRepository(opts ...func(*Repository)) *Repository {
	repo := &Repository{
		ID:       -1,
		Name:     "",
		UserID:   -1,
		Branches: []*Branch{},
	}

	for _, opt := range opts {
		opt(repo)
	}

	return repo
}

type Branch struct {
	ID   int64      `json:"id" db:"id"`
	Name string     `json:"name" db:"name"`
	Head *GardenTag `json:"head" db:"head"`
}

func NewBranch(opts ...func(*Branch)) *Branch {
	branch := &Branch{
		ID:   -1,
		Name: "",
		Head: NewGardenTag(),
	}

	for _, opt := range opts {
		opt(branch)
	}

	return branch
}

type GardenTag struct {
	ID        int64      `json:"id" db:"id"`
	Parent    *GardenTag `json:"parent" db:"parent"`
	Name      string     `json:"name" db:"name"`
	Signature string     `json:"signature" db:"signature"`
	Message   string     `json:"message" db:"message"`
	Timestamp time.Time  `json:"timestamp" db:"timestamp"`
	Tree      HashTree   `json:"tree" db:"tree"`
}

func NewGardenTag(opts ...func(*GardenTag)) *GardenTag {
	tag := &GardenTag{
		ID:        -1,
		Parent:    nil,
		Name:      "",
		Signature: "",
		Message:   "",
		Timestamp: time.Now(),
		Tree: HashTree{
			FolderNode: NewFolderNode(),
		},
	}

	for _, opt := range opts {
		opt(tag)
	}

	return tag
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

func NewFolderNode(opts ...func(*FolderNode)) *FolderNode {
	folder := &FolderNode{
		ID:        -1,
		Signature: "",
		Filename:  "",
		Path:      "",
		Contents: &struct {
			SubFolders []*FolderNode `json:"subfolders" db:"subfolders"`
			SubFiles   []*FileNode   `json:"subfiles" db:"subfiles"`
		}{
			SubFolders: []*FolderNode{},
			SubFiles:   []*FileNode{},
		},
	}

	for _, opt := range opts {
		opt(folder)
	}

	return folder
}

type FileNode struct {
	ID        int64  `json:"id" db:"id"`
	Signature string `json:"signature" db:"signature"`
	Filename  string `json:"filename" db:"filename"`
	Path      string `json:"path" db:"path"`
	Content   string `json:"content" db:"content"`
}

func NewFileNode(opts ...func(*FileNode)) *FileNode {
	file := &FileNode{
		ID:        -1,
		Signature: "",
		Filename:  "",
		Path:      "",
		Content:   "",
	}

	for _, opt := range opts {
		opt(file)
	}

	return file
}

type HashNode interface {
	FileNode | FolderNode
}
type HashTree struct {
	*FolderNode `json:"folder_node" db:"folder_node"`
}
