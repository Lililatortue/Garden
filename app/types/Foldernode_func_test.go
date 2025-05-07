package types

import (
	"testing"
)

func TestNewFolder(t *testing.T) {
	t.Logf("Testing Default NewFolder")
	t.Run("default", func(t *testing.T) {
		folder := NewFolderNode()
		if folder == nil {
			t.Errorf("expected folder to be non-nil")
		}
		if folder.ID != DefaultFolderNode.ID {
			t.Errorf("expected folder id to be %d, got %d", DefaultFolderNode.ID, folder.ID)
		}
		if folder.Signature != DefaultFolderNode.Signature {
			t.Errorf("expected folder signature to be empty, got %s", folder.Signature)
		}
		if folder.Filename != DefaultFolderNode.Filename {
			t.Errorf("expected folder filename to be empty, got %s", folder.Filename)
		}
		if folder.Path != DefaultFolderNode.Path {
			t.Errorf("expected folder path to be empty, got %s", folder.Path)
		}
	})
}

func FuzzNewFolderNode(f *testing.F) {
	f.Logf("Fuzzing NewFolderNode")
	testSet := []struct {
		id        int64
		signature string
		filename  string
		path      string
	}{
		{1, "test", "test", "test"},
		{2, "test2", "test2", "test2"},
		{3, "test3", "test3", "test3"},
		{4, "test4", "test4", "test4"},
	}
	for _, ts := range testSet {
		f.Add(ts.id, ts.signature, ts.filename, ts.path)
	}

	f.Fuzz(func(t *testing.T, id int64, signature string, filename string, path string) {
		t.Run("Test NewFolderNode with one function", func(t *testing.T) {
			folder := NewFolderNode(func(node *FolderNode) {
				node.ID = id
				node.Signature = signature
				node.Filename = filename
				node.Path = path
			})

			if folder == nil {
				t.Errorf("expected folder to be non-nil")
			}
			if folder.ID != id {
				t.Errorf("expected folder id to be %d, got %d", id, folder.ID)
			}
			if folder.Signature != signature {
				t.Errorf("expected folder signature to be %s, got %s", signature, folder.Signature)
			}
			if folder.Filename != filename {
				t.Errorf("expected folder filename to be %s, got %s", filename, folder.Filename)
			}
			if folder.Path != path {
				t.Errorf("expected folder path to be %s, got %s", path, folder.Path)
			}
		})

		t.Run("Test NewFolderNode with multiple functions", func(t *testing.T) {
			folder := NewFolderNode(func(node *FolderNode) {
				node.ID = id
			}, func(node *FolderNode) {
				node.Signature = signature
			}, func(node *FolderNode) {
				node.Filename = filename
			}, func(node *FolderNode) {
				node.Path = path
			})

			if folder == nil {
				t.Errorf("expected folder to be non-nil")
			}
			if folder.ID != id {
				t.Errorf("expected folder id to be %d, got %d", id, folder.ID)
			}
			if folder.Signature != signature {
				t.Errorf("expected folder signature to be %s, got %s", signature, folder.Signature)
			}
			if folder.Filename != filename {
				t.Errorf("expected folder filename to be %s, got %s", filename, folder.Filename)
			}
			if folder.Path != path {
				t.Errorf("expected folder path to be %s, got %s", path, folder.Path)
			}
		})
	})
}

func FuzzFolderNode_Traverse(f *testing.F) {
	f.Add(int64(1), int64(99))

	f.Fuzz(func(t *testing.T, firstID int64, secondID int64) {
		folder := NewFolderNode(func(node *FolderNode) {
			node.ID = firstID
			for range 10 {
				node.AddFolderNode(NewFolderNode(func(node *FolderNode) {
					node.ID = firstID
				}))
			}
		})

		folder.Traverse(func(node *FolderNode) {
			node.ID = secondID
		})

		folder.Traverse(func(node *FolderNode) {
			if node.ID != secondID {
				t.Errorf("expected node %d, got %d", secondID, node.ID)
			}
		})
	})
}

func FuzzFolderNode_TraverseAsync(f *testing.F) {
	f.Add(int64(1), int64(99))

	f.Fuzz(func(t *testing.T, firstID int64, secondID int64) {
		folder := NewFolderNode(func(node *FolderNode) {
			node.ID = firstID
			for range 10 {
				node.AddFolderNode(NewFolderNode(func(node *FolderNode) {
					node.ID = firstID
				}))
			}
		})

		folder.TraverseAsync(func(node *FolderNode) {
			node.ID = secondID
		})

		folder.Traverse(func(node *FolderNode) {
			if node.ID != secondID {
				t.Errorf("expected node %d, got %d", secondID, node.ID)
			}
		})
	})
}

func FuzzFolderNode_TraverseWithCondition(f *testing.F) {
	testSet := []struct {
		in int64
	}{
		{1},
		{2},
		{3},
		{4},
		{5},
		{6},
		{7},
		{8},
		{9},
	}

	for _, ts := range testSet {
		f.Add(ts.in)
	}

	f.Fuzz(func(t *testing.T, in int64) {
		folder := NewFolderNode(func(node *FolderNode) {
			node.ID = in
			node.Filename = "buzz"
			for range 10 {
				node.AddFolderNode(NewFolderNode(func(node *FolderNode) {
					node.ID = in
					node.Filename = "buzz"
				}))
			}
		})

		folder.TraverseWithCondition(func(n *FolderNode) {
			t.Logf("node %d has fizz filename", n.ID)
			n.Filename = "fizz"
		}, func(n *FolderNode) bool {
			t.Logf("checking node %d", n.ID)
			return (n.ID % 2) == 0
		})

		folder.Traverse(func(n *FolderNode) {
			t.Logf("node %d has fizz filename", n.ID)
			if n.ID != in {
				t.Errorf("expected node %d, got %d", in, n.ID)
			}
			if n.ID%2 == 0 && n.Filename != "fizz" {
				t.Errorf("expected node %d to have fizz filename, got %s", n.ID, n.Filename)
			} else if n.ID%2 != 0 && n.Filename != "buzz" {
				t.Errorf("expected node %d to have buzz filename, got %s", n.ID, n.Filename)
			}
		})
	})
}

func FuzzFolderNode_GetAllFileNodes(f *testing.F) {
	f.Logf("Fuzzing GetAllFileNodes")
	testSet := []struct {
		in int64
	}{
		{1},
		{12},
		{23},
		{34},
		{45},
		{56},
		{67},
		{78},
		{89},
		{91},
	}

	for _, ts := range testSet {
		f.Add(ts.in)
	}

	f.Fuzz(func(t *testing.T, in int64) {
		if in < 0 {
			return
		}

		folder := NewFolderNode(func(node *FolderNode) {
			node.ID = in
			node.AddFolderNode(NewFolderNode(func(node *FolderNode) {
				node.ID = in
			}))
			node.AddFileNode(NewFileNode())

			curr := node
			for range in {
				curr.AddFolderNode(NewFolderNode(func(node *FolderNode) {
					node.ID = in
				}))
				curr.AddFileNode(NewFileNode())
				curr = curr.SubFolders[0]
			}
		})

		nodes := folder.GetAllFileNodes()
		if len(nodes) != int(in)+1 {
			t.Errorf("expected %d nodes, got %d", in, len(nodes))
		}
	})
}

func FuzzFolderNode_GetAllFolderNodes(f *testing.F) {
	f.Logf("Fuzzing GetAllFolderNodes")
	testSet := []struct {
		in int64
	}{
		{1},
		{12},
		{23},
		{34},
		{45},
		{56},
		{67},
		{78},
		{89},
		{91},
	}

	for _, ts := range testSet {
		f.Add(ts.in)
	}

	f.Fuzz(func(t *testing.T, in int64) {
		if in < 0 {
			return
		}

		folder := NewFolderNode(func(node *FolderNode) {
			node.ID = in

			curr := node
			for range in {
				curr.AddFolderNode(NewFolderNode())
				curr = curr.SubFolders[0]
			}
		})

		nodes := folder.GetAllFolderNodes()
		if len(nodes) != int(in)+1 {
			t.Errorf("expected %d nodes, got %d", in, len(nodes)+1)
		}
	})
}

func FuzzFolderNode_AddFolderNode(f *testing.F) {
	f.Logf("Fuzzing AddFolderNode")
	testSet := []struct {
		in int64
	}{
		{1},
		{12},
		{23},
	}
	for _, ts := range testSet {
		f.Add(ts.in)
	}

	f.Fuzz(func(t *testing.T, in int64) {
		if in < 0 {
			return
		}

		folder := NewFolderNode(func(node *FolderNode) {
			node.ID = in

			for range in {
				node.AddFolderNode(NewFolderNode())
			}
		})
		nodes := folder.GetAllFolderNodes()
		if len(nodes) != int(in)+1 {
			t.Errorf("expected %d nodes, got %d", in, len(nodes)+1)
		}
	})
}

func FuzzFolderNode_AddFileNode(f *testing.F) {
	f.Logf("Fuzzing AddFileNode")
	testSet := []struct {
		in int64
	}{
		{1},
		{12},
		{23},
		{34},
		{45},
		{56},
		{67},
		{78},
		{89},
		{91},
	}
	for _, ts := range testSet {
		f.Add(ts.in)
	}

	f.Fuzz(func(t *testing.T, in int64) {
		t.Logf("in: %d", in)
		if in < 0 {
			return
		}
		folder := NewFolderNode(func(node *FolderNode) {
			node.ID = in
			for range in {
				node.AddFileNode(NewFileNode())
			}
		})

		nodes := folder.GetAllFileNodes()
		if len(nodes) != int(in) {
			t.Errorf("expected %d nodes, got %d", in, len(nodes))
		}
	})
}
