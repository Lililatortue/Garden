package hashtree

type FolderNode struct {
	Signature string
	Filename  string
	Path      string
	Contents  *struct {
		SubFolders []*FolderNode
		SubFiles   []*FileNode
	}
}

type FileNode struct {
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
