package data

import (
	"fmt"
	"garden/types"
	"log"
)

func (repo *GardenService) ReadUserByEmail(email string) (*types.User, error) {
	user, err := repo.Access.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}
	user.Repositories, err = repo.ReadRepositoryBy(int64(user.ID))
	if err != nil {
		return nil, fmt.Errorf("error reading repositories for user %s: %w", user.Name, err)
	}
	return user, nil
}

func (repo *GardenService) ReadUserByUsername(username string) (*types.User, error) {
	user, err := repo.Access.GetUserByUsername(username)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}
	user.Repositories, err = repo.ReadRepositoryBy(user.ID)
	if err != nil {
		return nil, fmt.Errorf("error reading repositories for user %s: %w", user.Name, err)
	}
	return user, nil
}

func (repo *GardenService) AddUser(user *types.User) (int64, error) {
	id, err := repo.Access.InsertUser(user)
	if err != nil {
		return -1, fmt.Errorf("error adding user: %w", err)
	}
	if user.Repositories != nil {
		for i, r := range user.Repositories {
			user.Repositories[i].ID, err = repo.AddRepository(r, id)
			if err != nil {
				return -1, fmt.Errorf("error adding repository %s: %w", r.Name, err)
			}

		}
	}

	return id, nil
}

func (repo *GardenService) ReadRepositoryBy(userId int64) ([]*types.Repository, error) {
	repos, err := repo.Access.GetRepositoriesForUser(userId)
	if err != nil {
		return nil, fmt.Errorf("error reading repositories: %w", err)
	}
	for _, code_repo := range repos {
		code_repo.Branches, err = repo.ReadBranchesBy(code_repo.ID)
		if err != nil {
			return nil, fmt.Errorf("error reading branches for repository %s: %w", code_repo.Name, err)
		}
	}

	return repos, nil
}

func (repo *GardenService) AddRepository(repository *types.Repository, userId int64) (int64, error) {
	repoId, err := repo.Access.InsertRepository(repository.Name, userId)
	if err != nil {
		return -1, fmt.Errorf("error adding repository: %w", err)
	}

	for i, branch := range repository.Branches {
		repository.Branches[i].ID, err = repo.AddBranch(branch, repoId)
		if err != nil {
			return 0, fmt.Errorf("error adding branch %s: %w", branch.Name, err)
		}
	}
	return repoId, nil
}

func (repo *GardenService) ReadBranchesBy(repoId int64) ([]*types.Branch, error) {
	branches, err := repo.Access.GetBranches(repoId)
	if err != nil {
		return nil, fmt.Errorf("error reading branches: %w", err)
	}

	for _, branch := range branches {
		branch.Head, err = repo.ReadTagBy(branch.Head.ID)
		if err != nil {
			return nil, fmt.Errorf("error reading head for branch %s: %w", branch.Name, err)
		}
	}
	return branches, nil
}

func (repo *GardenService) AddBranch(branch *types.Branch, repoId int64) (int64, error) {
	_, err := repo.AddTag(branch.Head)
	if err != nil {
		return 0, fmt.Errorf("error adding head for branch %s: %w", branch.Name, err)
	}
	return repo.Access.InsertBranch(branch, repoId)
}

func (repo *GardenService) UpdateBranchHead(branch *types.Branch) error {
	return repo.Access.UpdateBranchHead(branch)
}

func (repo *GardenService) ReadTagBy(tagId int64) (*types.GardenTag, error) {
	tag, err := repo.Access.GetGardenTag(tagId)
	if err != nil {
		return nil, fmt.Errorf("error reading tag: %w", err)
	}
	if tag.Parent != nil {
		tag.Parent, err = repo.ReadTagBy(tag.Parent.ID)
		if err != nil {
			return nil, fmt.Errorf("error reading parent for tag %s: %w", tag.Signature, err)
		}
	}

	tag.Tree, err = repo.ReadTree(tag.Tree.ID)
	if err != nil {
		return nil, fmt.Errorf("error reading tree for tag %s: %w", tag.Signature, err)
	}

	return tag, nil
}

func (repo *GardenService) AddTag(tag *types.GardenTag) (int64, error) {
	id, err := repo.Access.InsertGardenTag(tag)
	if err != nil {
		return -1, fmt.Errorf("error adding tag: %w", err)
	}

	if tag.Parent != nil {
		tag.Parent.ID, err = repo.AddTag(tag.Parent)
		if err != nil {
			return -1, fmt.Errorf("error adding parent for tag %s: %w", tag.Signature, err)
		}
	}

	return id, nil
}

func (repo *GardenService) ReadTree(treeId int64) (*types.HashTree, error) {
	folder, err := repo.ReadFolder(treeId)
	if err != nil {
		return nil, fmt.Errorf("error reading tree: %w", err)
	}

	return &types.HashTree{
		FolderNode: folder,
	}, nil
}

func (repo *GardenService) AddTree(tree *types.HashTree) (int64, error) {
	id, err := repo.AddFolder(tree.FolderNode, nil)
	if err != nil {
		return -1, fmt.Errorf("error adding tree: %w", err)
	}

	tree.TraverseAll(func(node *types.FolderNode) {
		for i, folder := range node.Contents.SubFolders {
			node.Contents.SubFolders[i].ID, err = repo.AddFolder(folder, &node.ID)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	})

	return id, nil
}

func (repo *GardenService) ReadFolder(folderId int64) (*types.FolderNode, error) {
	folder, err := repo.Access.GetFolder(folderId)
	if err != nil {
		return nil, fmt.Errorf("error reading folder: %w", err)
	}
	folder.Contents.SubFolders, err = repo.Access.GetTree(folder.ID)
	if err != nil {
		return nil, fmt.Errorf("error reading subfolders for folder %s: %w", folder.Filename, err)
	}

	for i, subFolder := range folder.Contents.SubFolders {
		folder.Contents.SubFolders[i], err = repo.ReadFolder(subFolder.ID)
		if err != nil {
			return nil, fmt.Errorf("error reading subfolders for folder %s: %w", folder.Filename, err)
		}
	}

	folder.Contents.SubFiles, err = repo.GetFilesFor(folderId)
	if err != nil {
		return nil, fmt.Errorf("error reading subfiles for folder %s: %w", folder.Filename, err)
	}

	return folder, nil
}

func (repo *GardenService) AddFolder(folder *types.FolderNode, parentID *int64) (int64, error) {
	id, err := repo.Access.InsertFolder(folder, parentID)
	if err != nil {
		return -1, fmt.Errorf("error adding folder: %w", err)
	}
	for i, file := range folder.Contents.SubFiles {
		folder.Contents.SubFiles[i].ID, err = repo.AddFile(file, id)
		if err != nil {
			log.Println(err.Error())
			return -1, fmt.Errorf("error adding file: %w", err)
		}
	}
	return id, nil
}

func (repo *GardenService) GetFilesFor(folderId int64) ([]*types.FileNode, error) {

	return repo.Access.GetFilesFor(folderId)
}

func (repo *GardenService) AddFile(file *types.FileNode, folderID int64) (int64, error) {
	return repo.Access.InsertFile(file, folderID)
}
