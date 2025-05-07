package sql

import (
	"fmt"
	"garden/types"
	"log"
	"time"
)

func (db *DBAccess) setup() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	db.mustCreateUserTable()
	db.mustCreateRepositoryTable()
	db.mustCreateFolderNodeTable()
	db.mustCreateGardenTagTable()
	db.mustCreateFileNodeTable()
	db.mustCreateBranchTable()

	// create default user and repository
	userId := db.mustSetupDefaultUser()
	db.mustSetupDefaultRepository(userId)

	log.Println("DB setup complete")
}

func (db *DBAccess) mustCreateGardenTagTable() {
	query := `CREATE TABLE IF NOT EXISTS GardenTag (
			id INTEGER PRIMARY KEY
				GENERATED ALWAYS AS IDENTITY,
			parent_id INTEGER,
			signature VARCHAR(40) NOT NULL,
			message TEXT,
			timestamp TIMESTAMP DEFAULT NOW() NOT NULL,
			tree_id INTEGER NOT NULL,
			FOREIGN KEY (parent_id) REFERENCES GardenTag (id)
			    ON DELETE CASCADE 
			    ON UPDATE CASCADE,
			FOREIGN KEY (tree_id) REFERENCES FolderNode (id)
				ON DELETE CASCADE 
				ON UPDATE CASCADE
				)`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating GardenTag table")
		fmt.Println(err.Error())
		panic(err)
	}
}

func (db *DBAccess) mustCreateFolderNodeTable() {
	query := `CREATE TABLE IF NOT EXISTS FolderNode (
			id INTEGER PRIMARY KEY
				GENERATED ALWAYS AS IDENTITY,  
			signature VARCHAR(40) NOT NULL,
			name VARCHAR(50) NOT NULL,
            parent_id INTEGER,
            FOREIGN KEY (parent_id) REFERENCES FolderNode (id)
			)`

	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating FolderNode table")
		fmt.Println(err.Error())
		panic(err)
	}
}

func (db *DBAccess) mustCreateBranchTable() {
	query := `CREATE TABLE IF NOT EXISTS Branch (
			id INTEGER PRIMARY KEY
				GENERATED ALWAYS AS IDENTITY,  
			name VARCHAR(50) NOT NULL,
			tag_id INTEGER NOT NULL,
			repository_id INTEGER NOT NULL,
			FOREIGN KEY (tag_id) REFERENCES GardenTag (id)
				ON DELETE CASCADE 
				ON UPDATE CASCADE,
			FOREIGN KEY (repository_id) REFERENCES Repository (id)
				ON DELETE CASCADE 
				ON UPDATE CASCADE,
			CONSTRAINT branch_name_repository_id_unique
				UNIQUE (name, repository_id)
			)`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating Branch table")
		fmt.Println(err.Error())
		panic(err)
	}
}

func (db *DBAccess) mustCreateFileNodeTable() {
	query := `CREATE TABLE IF NOT EXISTS FileNode (
			id INTEGER PRIMARY KEY
				GENERATED ALWAYS AS IDENTITY,  
			signature VARCHAR(40) NOT NULL,
			name VARCHAR(50) NOT NULL,
			content TEXT NOT NULL,
			folder_id INTEGER NOT NULL,
			FOREIGN KEY (folder_id) REFERENCES FolderNode (id)
				ON DELETE CASCADE 
				ON UPDATE CASCADE
			)`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating FileNode table")
		fmt.Println(err.Error())
		panic(err)
	}
}

func (db *DBAccess) mustCreateUserTable() {
	query := `CREATE TABLE IF NOT EXISTS GardenUser (
			id INTEGER PRIMARY KEY
				GENERATED ALWAYS AS IDENTITY,
			username VARCHAR(40) NOT NULL UNIQUE,
			password VARCHAR(40) NOT NULL,
			email VARCHAR(40) NOT NULL UNIQUE
			)`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating User table")
		fmt.Println(err.Error())
		panic(err)
	}
}

func (db *DBAccess) mustCreateRepositoryTable() {
	query := `CREATE TABLE IF NOT EXISTS Repository (
			id INTEGER PRIMARY KEY
				GENERATED ALWAYS AS IDENTITY,
			name VARCHAR(40) NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES GardenUser (id)
				ON DELETE CASCADE 
				ON UPDATE CASCADE,
			CONSTRAINT repository_name_user_id_unique
				UNIQUE (name, user_id)
			)`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating Repository table")
		fmt.Println(err.Error())
		panic(err)
	}
}

func (db *DBAccess) mustSetupDefaultUser() int64 {
	var (
		query = `
			INSERT INTO GardenUser (username, password, email) 
			VALUES ($1, $2, $3)
			ON CONFLICT (username) DO NOTHING
			RETURNING id`
		user = types.User{
			Name:     "test",
			Password: "test",
			Email:    "test@email.com",
		}
	)

	stmt, err := db.Prepare(query)
	if err != nil {
		panic(fmt.Errorf("error preparing statement: %w", err))
	}

	row := stmt.QueryRow(user.Name, user.Password, user.Email)

	var id int64
	if err = row.Scan(&id); err != nil {
		panic(fmt.Errorf("error inserting default user: %w", err))
	}

	return id

}

func (db *DBAccess) mustSetupDefaultRepository(userId int64) {
	var (
		repo = types.NewRepository(func(repository *types.Repository) {
			repository.Name = "test"
			repository.UserID = userId
			repository.Branches.Push(
				types.NewBranch(func(branch *types.Branch) {
					branch.Name = "main"
					branch.Head = *types.NewGardenTag(func(tag *types.GardenTag) {
						tag.Signature = "0000000000000000000000000000000000000000"
						tag.Message = "Initial commit"
						tag.Timestamp = time.Now()
						tag.Tree = types.HashTree{
							FolderNode: *types.NewFolderNode(func(folder *types.FolderNode) {
								folder.Filename = "test"
								folder.Signature = "0000000000000000000000000000000000000000"
								folder.Path = "/"
								folder.SubFiles.Push(
									types.NewFileNode(func(file *types.FileNode) {
										file.Filename = "test1"
										file.Signature = "0000000000000000000000000000000000000000"
										file.Content = "test1"
									}),
									types.NewFileNode(func(file *types.FileNode) {
										file.Filename = "test2"
										file.Signature = "0000000000000000000000000000000000000000"
										file.Content = "test2"
									}),
									types.NewFileNode(func(file *types.FileNode) {
										file.Filename = "test3"
										file.Signature = "0000000000000000000000000000000000000000"
										file.Content = "test3"
									}),
								)
								folder.SubFolders.Push(
									types.NewFolderNode(func(folder *types.FolderNode) {
										folder.Filename = "test4"
										folder.Signature = "0000000000000000000000000000000000000000"
										folder.Path = "/test4"
									}),
									types.NewFolderNode(func(folder *types.FolderNode) {
										folder.Filename = "test5"
										folder.Signature = "0000000000000000000000000000000000000000"
										folder.Path = "/test5"
										folder.SubFiles.Push(
											types.NewFileNode(func(file *types.FileNode) {
												file.Filename = "test6"
												file.Signature = "0000000000000000000000000000000000000000"
												file.Content = "test6"
											}),
											types.NewFileNode(func(file *types.FileNode) {
												file.Filename = "test7"
												file.Signature = "0000000000000000000000000000000000000000"
												file.Content = "test7"
											}),
											types.NewFileNode(func(file *types.FileNode) {
												file.Filename = "test8"
												file.Signature = "0000000000000000000000000000000000000000"
												file.Content = "test8"
											}),
										)
									}),
								)
							}),
						}
					})
				}),
			)
		})
	)

	if _, err := db.InsertRepository(repo.Name, repo.UserID); err != nil {
		panic(fmt.Errorf("error inserting default repository: %w", err))
	}

	for _, branch := range repo.Branches {
		if _, err := db.InsertBranch(branch, repo.ID); err != nil {
			panic(fmt.Errorf("error inserting default branch: %w", err))
		}

		for tag := range branch.Head.IterateToParent() {
			parentId, err := db.InsertFolder(&tag.Tree.FolderNode, nil)
			if err != nil {
				panic(fmt.Errorf("error inserting default folder: %w", err))
			}
			tag.Tree.ID = parentId

			tag.Tree.Traverse(func(node *types.FolderNode) {
				for _, file := range node.SubFiles {
					id, err := db.InsertFile(file, node.ID)
					if err != nil {
						panic(fmt.Errorf("error inserting default file: %w", err))
					}
					file.ID = id
				}

				for _, folder := range node.SubFolders {
					id, err := db.InsertFolder(folder, &node.ID)
					if err != nil {
						panic(fmt.Errorf("error inserting default folder: %w", err))
					}
					folder.ID = id
				}
			})
		}
	}

}
