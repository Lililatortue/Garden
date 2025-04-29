package data

import (
	"fmt"
	"garden/types"
	"log"
	"time"
)

func (gs *GardenService) ReadUserByEmail(email string) (*types.User, error) {
	user, err := gs.Access.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}
	user.Repositories, err = gs.ReadRepositoryBy(int64(user.ID))
	if err != nil {
		return nil, fmt.Errorf("error reading repositories for user %s: %w", user.Name, err)
	}
	return user, nil
}

func (gs *GardenService) ReadUserByUsername(username string) (*types.User, error) {
	user, err := gs.Access.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}
	user.Repositories, err = gs.ReadRepositoryBy(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error reading repositories for user %s: %w", user.Name, err)
	}
	return user, nil
}

func (gs *GardenService) AddUser(user *types.User) (int64, error) {
	id, err := gs.Access.InsertUser(user)
	if err != nil {
		return -1, fmt.Errorf("error adding user: %w", err)
	}
	if user.Repositories != nil {
		for i, r := range user.Repositories {
			user.Repositories[i].ID, err = gs.AddRepository(r, id)
			if err != nil {
				return -1, fmt.Errorf("error adding repository %s: %w", r.Name, err)
			}

		}
	}

	return id, nil
}

func (gs *GardenService) ReadRepositoryBy(userId int64) ([]*types.Repository, error) {
	log.Printf("Reading repositories for user %d\n", userId)
	repos, err := gs.Access.GetRepositoriesForUser(userId)
	if err != nil {
		return nil, fmt.Errorf("error reading repositories: %w", err)
	}
	log.Printf("Read %d repositories\n", len(repos))
	for i, repo := range repos {

		repos[i].Branches, err = gs.ReadBranchesOf(repo.ID)
		if err != nil {
			return nil, fmt.Errorf("error reading branches for repository %s: %w", repo.Name, err)
		}
	}

	return repos, nil
}

func (gs *GardenService) ReadRepository(name string, userID int64) (*types.Repository, error) {
	repo, err := gs.Access.GetRepositoryByName(name, userID)
	if err != nil {
		return nil, fmt.Errorf("error reading repository: %w", err)
	}
	return repo, nil
}

func (gs *GardenService) AddRepository(repository *types.Repository, userId int64) (int64, error) {
	repoId, err := gs.Access.InsertRepository(repository.Name, userId)
	if err != nil {
		return -1, fmt.Errorf("error adding repository: %w", err)
	}

	repository.Branches, err = gs.ReadBranchesOf(repoId)

	for i, branch := range repository.Branches {
		repository.Branches[i].ID, err = gs.AddBranch(branch, repoId)
		if err != nil {
			return 0, fmt.Errorf("error adding branch %s: %w", branch.Name, err)
		}
	}
	return repoId, nil
}

func (gs *GardenService) InitRepository(repoName string, userId int64) (*types.Repository, error) {
	var (
		repo = &types.Repository{
			Name: repoName,
		}
		branch = &types.Branch{
			Name: "main",
		}
		tag = &types.GardenTag{
			Name:      "main",
			Signature: "b28b7af69320201d1cf206ebf28373980add1451",
			Parent:    nil,
			Message:   "Initial commit",
			Timestamp: time.Now(),
		}
		tree = &types.HashTree{
			FolderNode: &types.FolderNode{
				Filename:  "test",
				Path:      "/",
				Signature: "42099b4af021e53fd8fd4e056c2568d7c2e3ffa8",
			},
		}
	)
	tag.Tree.FolderNode.Contents.SubFiles = []*types.FileNode{
		{
			Filename:  "README.md",
			Path:      "/README.md",
			Signature: "e1d57665c76144e7bb6a1436c4be9213d2610534",
			Content:   "# test\n",
		},
	}

	repoId, err := gs.AddRepository(repo, userId)
	if err != nil {
		return nil, fmt.Errorf("error adding repository: %w", err)
	}
	repo.ID = repoId

	branch.ID, err = gs.AddBranch(&types.Branch{Name: "main"}, repoId)
	if err != nil {
		return nil, fmt.Errorf("error adding branch: %w", err)
	}
	repo.Branches = append(repo.Branches, branch)

	tag.ID, err = gs.AddTag(tag)
	if err != nil {
		return nil, fmt.Errorf("error adding tag: %w", err)
	}

	tree.ID, err = gs.AddTree(tree)
	if err != nil {
		return nil, fmt.Errorf("error adding tree: %w", err)
	}
	tag.Tree = *tree

	return repo, nil

}

func (gs GardenService) ReadBranch(branchName string, repoId int64) (*types.Branch, error) {
	branch, err := gs.Access.GetBranch(branchName, repoId)
	if err != nil {
		return nil, fmt.Errorf("error reading branch: %w", err)
	}

	return branch, nil

}

func (gs *GardenService) ReadBranchesOf(repoId int64) ([]*types.Branch, error) {
	branches, err := gs.Access.GetBranches(repoId)
	if err != nil {
		return nil, fmt.Errorf("error reading branches: %w", err)
	}

	for i, branch := range branches {
		branches[i].Head, err = gs.ReadTagRecursiveBy(branch.Head.ID)
		if err != nil {
			return nil, fmt.Errorf("error reading head for branch %s: %w", branch.Name, err)
		}
	}
	return branches, nil
}

func (gs *GardenService) AddBranch(branch *types.Branch, repoId int64) (int64, error) {
	_, err := gs.AddTag(branch.Head)
	if err != nil {
		return 0, fmt.Errorf("error adding head for branch %s: %w", branch.Name, err)
	}
	return gs.Access.InsertBranch(branch, repoId)
}

func (gs *GardenService) UpdateBranchHead(branch *types.Branch) error {
	return gs.Access.UpdateBranchHead(branch)
}

func (gs GardenService) ReadTagRecursiveBy(tagId int64) (*types.GardenTag, error) {
	tag := types.NewGardenTag(func(tag *types.GardenTag) {
		tag.ID = tagId
	})

	for currTag := range tag.IterateToParent() {
		currTag, err := gs.ReadTagBy(currTag.ID)
		if err != nil {
			return nil, fmt.Errorf("error reading tag: %w", err)
		}
		tree, err := gs.ReadTree(currTag.Tree.ID)
		if err != nil {
			return nil, fmt.Errorf("error reading tree for tag %s: %w", currTag.Signature, err)
		}
		currTag.Tree = *tree
	}
	return tag, nil
}

func (gs *GardenService) ReadTagBy(tagId int64) (*types.GardenTag, error) {
	tag, err := gs.Access.GetGardenTag(tagId)
	if err != nil {
		return nil, fmt.Errorf("error reading tag: %w", err)
	}

	tree, err := gs.ReadTree(tag.Tree.ID)
	if err != nil {
		return nil, fmt.Errorf("error reading tree for tag %s: %w", tag.Signature, err)
	}

	tag.Tree = *tree

	return tag, nil
}

func (gs *GardenService) AddTag(tag *types.GardenTag) (int64, error) {
	id, err := gs.Access.InsertGardenTag(tag)
	if err != nil {
		return -1, fmt.Errorf("error adding tag: %w", err)
	}

	if tag.Parent != nil {
		tag.Parent.ID, err = gs.AddTag(tag.Parent)
		if err != nil {
			return -1, fmt.Errorf("error adding parent for tag %s: %w", tag.Signature, err)
		}
	}

	return id, nil
}

func (gs *GardenService) ReadTree(treeId int64) (*types.HashTree, error) {
	var (
		tree = types.HashTree{}
	)
	root, err := gs.Access.GetFolder(treeId)
	if err != nil {
		return nil, fmt.Errorf("error reading tree: %w", err)
	}
	tree.FolderNode = root

	tree.Traverse(func(node *types.FolderNode) {
		for i, folder := range node.Contents.SubFolders {
			readFolder, err := gs.ReadFolder(folder.ID)
			if err != nil {
				log.Println(err.Error())
			}
			node.Contents.SubFolders[i] = readFolder
		}
	})
	return &types.HashTree{FolderNode: root}, nil
}

func (gs *GardenService) AddTree(tree *types.HashTree) (int64, error) {
	id, err := gs.AddFolder(tree.FolderNode, nil)
	if err != nil {
		return -1, fmt.Errorf("error adding tree: %w", err)
	}

	tree.Traverse(func(node *types.FolderNode) {
		for i, folder := range node.Contents.SubFolders {
			node.Contents.SubFolders[i].ID, err = gs.AddFolder(folder, &node.ID)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	})

	return id, nil
}

func (gs *GardenService) ReadFolder(folderId int64) (*types.FolderNode, error) {
	folder, err := gs.Access.GetFolder(folderId)
	if err != nil {
		return nil, fmt.Errorf("error reading folder: %w", err)
	}
	folder.Contents.SubFolders, err = gs.Access.GetSubFolders(folder.ID)
	if err != nil {
		return nil, fmt.Errorf("error reading subfolders for folder %s: %w", folder.Filename, err)
	}

	folder.Contents.SubFiles, err = gs.GetFilesFor(folderId)
	if err != nil {
		return nil, fmt.Errorf("error reading subfiles for folder %s: %w", folder.Filename, err)
	}

	return folder, nil
}

func (gs *GardenService) AddFolder(folder *types.FolderNode, parentID *int64) (int64, error) {
	id, err := gs.Access.InsertFolder(folder, parentID)
	if err != nil {
		return -1, fmt.Errorf("error adding folder: %w", err)
	}
	for i, file := range folder.Contents.SubFiles {
		folder.Contents.SubFiles[i].ID, err = gs.AddFile(file, id)
		if err != nil {
			log.Println(err.Error())
			return -1, fmt.Errorf("error adding file: %w", err)
		}
	}
	return id, nil
}

func (gs *GardenService) GetFilesFor(folderId int64) ([]*types.FileNode, error) {

	return gs.Access.GetFilesFor(folderId)
}

func (gs *GardenService) AddFile(file *types.FileNode, folderID int64) (int64, error) {
	return gs.Access.InsertFile(file, folderID)
}
