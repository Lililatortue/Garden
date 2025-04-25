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

func (f *FolderNode) TraverseAll(action func(node *FolderNode)) {
	var waitGrp sync.WaitGroup
	for n := range f.Iterate() {
		waitGrp.Add(1)
		go n.TraverseAll(action)
		action(n)
		waitGrp.Done()
	}
	waitGrp.Wait()
}

func (f *FolderNode) Traverse(action func(n *FolderNode), predicate func(n *FolderNode) bool) {
	var waitGrp sync.WaitGroup
	for n := range f.Iterate() {
		waitGrp.Add(1)
		go n.Traverse(action, predicate)
		action(n)
		waitGrp.Done()
	}
	waitGrp.Wait()
}

func (f *FolderNode) GetAllFileNodes() []*FileNode {
	var nodes []*FileNode
	f.TraverseAll(func(n *FolderNode) {
		nodes = append(nodes, n.Contents.SubFiles...)
	})
	return nodes
}
func (f *FolderNode) GetAllFolderNodes() []*FolderNode {
	var nodes []*FolderNode
	f.TraverseAll(func(n *FolderNode) {
		nodes = append(nodes, n.Contents.SubFolders...)
	})
	return nodes
}

func (f *FolderNode) GetSubFolders(predicate func(node *FolderNode) bool) []*FolderNode {
	var nodes []*FolderNode
	f.Traverse(func(n *FolderNode) {
		nodes = append(nodes, n.Contents.SubFolders...)
	}, predicate)
	return nodes
}

func (f *FolderNode) GetFiles(predicate func(node *FolderNode) bool) []*FileNode {
	var nodes []*FileNode
	f.Traverse(func(n *FolderNode) {
		nodes = append(nodes, n.Contents.SubFiles...)
	}, predicate)
	return nodes
}

func (f *FolderNode) AddFolderNode(node FolderNode) {
	f.Contents.SubFolders = append(f.Contents.SubFolders, &node)
}

func (f *FolderNode) AddFileNode(node FileNode) {
	f.Contents.SubFiles = append(f.Contents.SubFiles, &node)
}
