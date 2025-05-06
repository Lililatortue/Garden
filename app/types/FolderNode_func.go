package types

import (
	"sync"
)

// NewFolderNode Creates a new folder node with default values adequate for a folder node.
// If options are provided they are applied to the node.
// The options are applied in the order they are provided.
func NewFolderNode(opts ...func(*FolderNode)) *FolderNode {
	folder := DefaultFolderNode

	for _, opt := range opts {
		opt(&folder)
	}

	return &folder
}

// Traverse Function that traverse each node recursively and apply the action to that node
func (f *FolderNode) Traverse(action func(node *FolderNode)) {
	action(f)
	for _, n := range f.SubFolders {
		n.Traverse(action)
	}
}

func (f *FolderNode) TraverseAsync(action func(node *FolderNode)) {
	waitGrp := &sync.WaitGroup{}
	action(f)
	for _, n := range f.SubFolders {
		waitGrp.Add(1)
		go n.traverseWithWaitGroup(action, waitGrp)
	}
	waitGrp.Wait()
}

func (f *FolderNode) traverseWithWaitGroup(action func(node *FolderNode), waitGrp *sync.WaitGroup) {
	action(f)
	for _, n := range f.SubFolders {
		waitGrp.Add(1)
		go n.traverseWithWaitGroup(action, waitGrp)
	}
	waitGrp.Done()
}

func (f *FolderNode) TraverseWithCondition(action func(n *FolderNode), predicate func(n *FolderNode) bool) {
	if !predicate(f) {
		return
	}
	action(f)
	for _, n := range f.SubFolders {
		n.TraverseWithCondition(action, predicate)
	}
}

func (f *FolderNode) GetAllFileNodes() []*FileNode {
	var nodes []*FileNode
	lock := &sync.Mutex{}
	f.TraverseAsync(func(n *FolderNode) {
		lock.Lock()
		nodes = append(nodes, n.SubFiles...)
		lock.Unlock()
	})
	return nodes
}

func (f *FolderNode) GetAllFolderNodes() []*FolderNode {
	var (
		nodes = []*FolderNode{
			f,
		}
	)
	lock := &sync.Mutex{}
	f.TraverseAsync(func(n *FolderNode) {
		lock.Lock()
		nodes = append(nodes, n.SubFolders...)
		lock.Unlock()
	})
	return nodes
}

func (f *FolderNode) AddFolderNode(node ...*FolderNode) {
	f.SubFolders.Push(node...)
}

func (f *FolderNode) AddFileNode(node ...*FileNode) {
	f.SubFiles.Push(node...)
}
