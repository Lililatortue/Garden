package sql

import (
	"garden/gardentag"
	"garden/hashtree"
)

func (db *DBAccess) insertGardenTag(tag gardentag.GardenTag, repoID int, treeId int) (int64, error) {
	stmt, err := db.Prepare(`
		INSERT INTO GardenTag (
			signature, 
			message,
			timestamp, 
			tree_id, 
			repository_id
		) VALUES ($1, $2, $3, $4)
		`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(tag.Signature, tag.Message, tag.Timestamp, treeId, repoID)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (db *DBAccess) insertRepository(repo_name string, userId int64) (int64, error) {
	stmt, err := db.Prepare(`
		INSERT INTO Repository (name, user_id)
		VALUES ($1, $2)
		`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(repo_name, userId)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (db *DBAccess) insertFolder(folder hashtree.FolderNode) (int64, error) {
	stmt, err := db.Prepare(`
		INSERT INTO FolderNode (signature, name)
		VALUES ($1, $2)
		`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(folder.Signature, folder.Filename)
	if err != nil {
		return -1, err
	}

	return res.LastInsertId()
}

func (db *DBAccess) insertFile(file hashtree.FileNode, folderId int64) (int64, error) {
	stmt, err := db.Prepare(`
		INSERT INTO FileNode 
		  (signature, name, content, folder_id)
		VALUES ($1, $2, $3, $4)
		`)
	if err != nil {
		return -1, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(file.Signature, file.Filename, file.Content, folderId)
	if err != nil {
		return -1, err
	}
	return res.LastInsertId()
}
