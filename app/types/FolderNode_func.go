package types

import (
	"iter"
	"sync"
)

type NodeIterator iter.Seq[*FolderNode]

func (f *FolderNode) Iterate() NodeIterator {
	return func(yield func(node *FolderNode) bool) {
		for _, n := range f.Contents.SubFolders {
			if !yield(n) {
				return
			}
		}
	}
}

func (f *FolderNode) Traverse(action func(node *FolderNode)) {
	action(f)
	for n := range f.Iterate() {
		n.Traverse(action)
	}
}
func (f *FolderNode) TraverseAsync(action func(node *FolderNode)) {
	waitGrp := &sync.WaitGroup{}
	action(f)
	for n := range f.Iterate() {
		waitGrp.Add(1)
		go n.traverseWithWaitGroup(action, waitGrp)
	}
	waitGrp.Wait()
}

func (f *FolderNode) traverseWithWaitGroup(action func(node *FolderNode), waitGrp *sync.WaitGroup) {
	action(f)
	for n := range f.Iterate() {
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
	for n := range f.Iterate() {
		n.TraverseWithCondition(action, predicate)
	}
}

func (f *FolderNode) GetAllFileNodes() []*FileNode {
	var nodes []*FileNode
	lock := &sync.Mutex{}
	f.TraverseAsync(func(n *FolderNode) {
		lock.Lock()
		nodes = append(nodes, n.Contents.SubFiles...)
		lock.Unlock()
	})
	return nodes
}
func (f *FolderNode) GetAllFolderNodes() []*FolderNode {
	var nodes []*FolderNode
	lock := &sync.Mutex{}
	f.TraverseAsync(func(n *FolderNode) {
		lock.Lock()
		nodes = append(nodes, n.Contents.SubFolders...)
		lock.Unlock()
	})
	return nodes
}

func (f *FolderNode) GetSubFolders(predicate func(node *FolderNode) bool) []*FolderNode {
	var nodes []*FolderNode
	f.TraverseWithCondition(func(n *FolderNode) {
		nodes = append(nodes, n.Contents.SubFolders...)
	}, predicate)
	return nodes
}

func (f *FolderNode) GetSubFiles(predicate func(node *FolderNode) bool) []*FileNode {
	var nodes []*FileNode
	f.TraverseWithCondition(func(n *FolderNode) {
		nodes = append(nodes, n.Contents.SubFiles...)
	}, predicate)
	return nodes
}

func (f *FolderNode) AddFolderNode(node *FolderNode) {
	f.Contents.SubFolders = append(f.Contents.SubFolders, node)
}

func (f *FolderNode) AddFileNode(node *FileNode) {
	f.Contents.SubFiles = append(f.Contents.SubFiles, node)
}
