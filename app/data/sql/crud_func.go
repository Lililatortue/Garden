package sql

import (
	"garden/types"
)

func (db *DBAccess) GetRepository(repoId int64) (*types.Repository, error) {
	stmt, err := db.Prepare(`
		SELECT * FROM Repository WHERE id = $1
		`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(repoId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var repo types.Repository
	for rows.Next() {
		err = rows.Scan(&repo.ID, &repo.Name, &repo.UserID)
		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &repo, nil
}

func (db *DBAccess) GetGardenTag(tagId int64) (*types.GardenTag, error) {
	var (
		tag    types.GardenTag
		repoId int64
	)

	stmt, err := db.Prepare(`
		SELECT * FROM GardenTag WHERE id = $1
		`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(tagId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&tag.ID,
			&tag.Signature,
			&tag.Message,
			&tag.Timestamp,
			&tag.Tree.ID,
			&repoId,
		)
		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &tag, nil
}

func (db *DBAccess) GetFolder(folderId int64) (*types.FolderNode, error) {
	stmt, err := db.Prepare(`
		SELECT * FROM FolderNode WHERE id = $1
		`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(folderId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folder types.FolderNode
	for rows.Next() {
		err = rows.Scan(
			&folder.ID,
			&folder.Signature,
			&folder.Filename,
		)
		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &folder, nil
}

func (db *DBAccess) GetFile(fileId int64) (*types.FileNode, error) {
	stmt, err := db.Prepare(`
		SELECT * FROM FileNode WHERE id = $1
		`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(fileId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var file types.FileNode
	var folderId int64
	for rows.Next() {
		err = rows.Scan(
			&file.ID,
			&file.Signature,
			&file.Filename,
			&file.Content,
			&folderId,
		)
		if err != nil {
			return nil, err
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &file, nil
}

func (db *DBAccess) InsertGardenTag(tag types.GardenTag, repoID int, treeId int) (int64, error) {
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

func (db *DBAccess) InsertRepository(repo_name string, userId int64) (int64, error) {
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

func (db *DBAccess) InsertFolder(folder types.FolderNode) (int64, error) {
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

func (db *DBAccess) InsertFile(file types.FileNode, folderId int64) (int64, error) {
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
