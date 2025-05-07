package types

import "time"

var (
	DefaultUser = User{
		ID:           -1,
		Name:         "",
		Email:        "",
		Password:     "",
		Repositories: []*Repository{},
	}

	DefaultRepository = Repository{
		ID:       -1,
		Name:     "",
		UserID:   -1,
		Branches: []*Branch{},
	}

	DefaultBranch = Branch{
		ID:   -1,
		Name: "",
		Head: DefaultGardenTag,
	}

	DefaultGardenTag = GardenTag{
		ID:        -1,
		Parent:    nil,
		Name:      "",
		Signature: "",
		Message:   "",
		Timestamp: time.UnixMicro(0),
		Tree: HashTree{
			FolderNode: DefaultFolderNode,
		},
	}

	DefaultFolderNode = FolderNode{
		ID:         -1,
		Signature:  "",
		Filename:   "",
		Path:       "",
		SubFolders: []*FolderNode{},
		SubFiles:   []*FileNode{},
	}

	DefaultFileNode = FileNode{
		ID:        -1,
		Signature: "",
		Filename:  "",
		Path:      "",
		Content:   "",
	}
)
