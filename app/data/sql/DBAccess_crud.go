package sql

import (
	"database/sql"
	"fmt"
	"garden/types"
	"log"
)

func (db *DBAccess) GetUserByEmail(email string) (*types.User, error) {
	var user types.User

	stmt, err := db.Prepare(`
		SELECT * FROM GardenUser WHERE email = $1
		`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	rows, err := stmt.Query(email)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(rows)

	if rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Password,
			&user.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
	}
	return &user, nil
}

func (db *DBAccess) GetUserByUsername(username string) (*types.User, error) {
	var user types.User
	stmt, err := db.Prepare(`SELECT * FROM GardenUser WHERE username = $1`)
	if err != nil {
		return nil, fmt.Errorf("1error preparing statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	rows, err := stmt.Query(username)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}

	if rows.Next() {
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Password,
			&user.Email,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning user: %w", err)
		}
	}

	return &user, nil
}

func (db *DBAccess) InsertUser(user *types.User) (int64, error) {
	stmt, err := db.Prepare(`
		INSERT INTO GardenUser (username, password, email)
		VALUES ($1, $2, $3)
		RETURNING id
		`)
	if err != nil {
		return -1, err
	}

	row := stmt.QueryRow(user.Name, user.Password, user.Email)
	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting user: %w", err)
	}
	return id, nil
}

func (db *DBAccess) GetRepository(repoId int64) (*types.Repository, error) {
	stmt, err := db.Prepare(`
		SELECT * FROM Repository WHERE id = $1
		`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	rows, err := stmt.Query(repoId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(rows)

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

func (db *DBAccess) GetRepositoryByName(name string, userID int64) (*types.Repository, error) {
	query := `SELECT * FROM Repository 
				WHERE name = $1 
		  		AND user_id = $2`

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}()

	row := stmt.QueryRow(name, userID)

	var repo types.Repository
	if err = row.Scan(&repo.ID, &repo.Name, &repo.UserID); err != nil {
		return nil, fmt.Errorf("error scanning repository: %w", err)
	}

	return &repo, nil
}

func (db *DBAccess) GetRepositoriesForUser(userId int64) ([]*types.Repository, error) {
	var repos []*types.Repository = make([]*types.Repository, 0)
	stmt, err := db.Prepare(`
		SELECT * FROM Repository WHERE user_id = $1
		`)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	rows, err := stmt.Query(userId)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}

	for rows.Next() {
		var repo types.Repository
		err := rows.Scan(
			&repo.ID,
			&repo.Name,
			&repo.UserID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning repository: %w", err)
		}
		repos = append(repos, &repo)
	}

	return repos, nil
}

func (db *DBAccess) InsertRepository(repoName string, userId int64) (int64, error) {
	stmt, err := db.Prepare(`
		INSERT INTO Repository (name, user_id)
		VALUES ($1, $2)
		RETURNING id
		`)
	if err != nil {
		return -1, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	row := stmt.QueryRow(repoName, userId)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error scanning repository id: %w", err)
	}

	return id, nil
}

func (db *DBAccess) GetBranch(name string, repoId int64) (*types.Branch, error) {
	var (
		query = `
			SELECT id, tag_id FROM Branch 
			  WHERE name = $1 
			  AND repository_id = $2`
		branch types.Branch = types.Branch{
			Name: name,
			Head: &types.GardenTag{},
		}
	)

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}

	row := stmt.QueryRow(name, repoId)

	err = row.Scan(
		&branch.ID,
		&branch.Head.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("error scanning branch: %w", err)
	}

	return &branch, nil
}

func (db *DBAccess) GetBranches(repoId int64) ([]*types.Branch, error) {
	var branches []*types.Branch = make([]*types.Branch, 0)
	stmt, err := db.Prepare(`
		SELECT * FROM Branch WHERE repository_id = $1
		`)
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	rows, err := stmt.Query(repoId)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}

	for rows.Next() {
		var branch types.Branch

		err := rows.Scan(
			&branch.ID,
			&branch.Name,
			&branch.Head.ID,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning branch: %w", err)
		}
		branches = append(branches, &branch)
	}

	return branches, nil

}

func (db *DBAccess) InsertBranch(branch *types.Branch, repoID int64) (int64, error) {
	stmt, err := db.Prepare(`
		Insert into Branch 
			(name, tag_id, repository_id)
		values ($1, $2, $3)
		RETURNING id
		`)
	if err != nil {
		return -1, fmt.Errorf("error inserting branch: %w", err)
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	row := stmt.QueryRow(branch.Name, branch.Head.ID, repoID)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error scanning branch id: %w", err)
	}

	return id, nil
}

func (db *DBAccess) UpsertBranch(branch *types.Branch, repoID int64) (int64, error) {
	stmt, err := db.Prepare(`
		INSERT INTO Branch (name, tag_id, repository_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (name, repository_id) DO UPDATE SET tag_id = $2
		RETURNING id
		`)
	if err != nil {
		return -1, fmt.Errorf("error upserting branch: %w", err)
	}

	row := stmt.QueryRow(branch.Name, branch.Head.ID, repoID)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error scanning branch id: %w", err)
	}

	return id, nil

}

func (db *DBAccess) UpdateBranchHead(branch *types.Branch) error {
	stmt, err := db.Prepare(`UPDATE Branch SET head_id = $1 WHERE id = $2`)
	if err != nil {
		return fmt.Errorf("Error formating query: %w", err)
	}

	_, err = stmt.Exec(branch.Head.ID, branch.ID)
	if err != nil {
		return fmt.Errorf("Error updating branch: %w", err)
	}

	return nil
}

func (db *DBAccess) GetGardenTag(tagId int64) (*types.GardenTag, error) {
	var (
		tag types.GardenTag
	)

	stmt, err := db.Prepare(`
		SELECT * FROM GardenTag WHERE id = $1
		`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	rows, err := stmt.Query(tagId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(rows)

	if rows.Next() {
		err = rows.Scan(
			&tag.ID,
			&tag.Parent.ID,
			&tag.Signature,
			&tag.Message,
			&tag.Timestamp,
			&tag.Tree.ID,
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

func (db *DBAccess) InsertGardenTag(tag *types.GardenTag) (int64, error) {
	var (
		query  string
		params []any
	)

	if tag.Parent == nil {
		query = `INSERT INTO GardenTag 
    				(signature, message, timestamp, tree_id) 
				 VALUES ($1, $2, $3, $4)
				 RETURNING id`
		params = []any{tag.Signature, tag.Message, tag.Timestamp, tag.Tree.ID}
	} else {
		query = `INSERT INTO GardenTag 
    				(parent_id, signature, message, timestamp, tree_id) 
				 VALUES ($1, $2, $3, $4, $5)
				 RETURNING id`
		params = []any{tag.Parent.ID, tag.Signature, tag.Message, tag.Timestamp, tag.Tree.ID}
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	row := stmt.QueryRow(params)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error scanning tag id: %w", err)
	}

	return id, nil
}

func (db *DBAccess) GetFolder(folderId int64) (*types.FolderNode, error) {
	stmt, err := db.Prepare(`
		SELECT * FROM FolderNode WHERE id = $1
		`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	rows, err := stmt.Query(folderId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(rows)

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

func (db *DBAccess) GetTree(treeId int64) ([]*types.FolderNode, error) {
	var folders []*types.FolderNode = make([]*types.FolderNode, 0)

	stmt, err := db.Prepare(`
		SELECT * FROM FolderNode WHERE parent_id = $1
		`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	rows, err := stmt.Query(treeId)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(rows)

	for rows.Next() {
		var folder types.FolderNode
		err := rows.Scan(
			&folder.ID,
			&folder.Signature,
			&folder.Filename,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning folder: %w", err)
		}
		folders = append(folders, &folder)
	}

	return folders, nil
}

func (db *DBAccess) InsertFolder(folder *types.FolderNode, parentID *int64) (int64, error) {
	var (
		query  string
		params = []interface{}{folder.Signature, folder.Filename}
	)
	if parentID == nil {
		query = `INSERT INTO FolderNode (signature, name) VALUES ($1, $2) RETURNING id`
	} else {
		query = `INSERT INTO FolderNode (signature, name, parent_id) VALUES ($1, $2, $3) RETURNING id`
		params = append(params, *parentID)

	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	row := stmt.QueryRow(params...)

	var id int64
	err = row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error scanning folder id: %w", err)
	}

	return id, nil
}

func (db *DBAccess) GetFilesFor(folderId int64) ([]*types.FileNode, error) {
	var files []*types.FileNode = make([]*types.FileNode, 0)

	stmt, err := db.Prepare(`SELECT * FROM FileNode WHERE folder_id = $1`)
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	rows, err := stmt.Query(folderId)
	if err != nil {
		return nil, fmt.Errorf("error executing statement: %w", err)
	}

	for rows.Next() {
		var file types.FileNode
		err := rows.Scan(&file.ID, &file.Signature, &file.Filename, &file.Content)

		if err != nil {
			return nil, fmt.Errorf("error scanning file: %w", err)
		}
		files = append(files, &file)
	}

	return files, nil
}

func (db *DBAccess) InsertFile(file *types.FileNode, folderId int64) (int64, error) {
	stmt, err := db.Prepare(`
		INSERT INTO FileNode 
		  (signature, name, content, folder_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id
		`)
	if err != nil {
		return -1, err
	}
	defer func(stmt *sql.Stmt) {
		err := stmt.Close()
		if err != nil {
			log.Println(err.Error())
		}
	}(stmt)

	row, err := stmt.Query(file.Signature, file.Filename, file.Content, folderId)
	if err != nil {
		return -1, err
	}
	var id int64
	if row.Next() {
		err = row.Scan(&id)
		if err != nil {
			return -1, fmt.Errorf("error scanning file id: %w", err)
		}
		return id, nil
	}

	return -1, ErrNoReturningID
}
