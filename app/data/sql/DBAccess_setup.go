package sql

import (
	"database/sql"
	"errors"
	"fmt"
	"garden/types"
	"log"
)

func (db *DBAccess) setup() {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()

	db.createUserTable()
	db.createRepositoryTable()
	db.createFolderNodeTable()
	db.createGardenTagTable()
	db.createFileNodeTable()
	db.createBranchTable()

	// create default user and repository
	userId := db.setupDefaultUser()
	_ = db.setupDefaultRepository(userId)

	log.Println("DB setup complete")
}

func (db *DBAccess) createGardenTagTable() {
	query := `CREATE TABLE IF NOT EXISTS GardenTag (
			id INTEGER PRIMARY KEY
				GENERATED ALWAYS AS IDENTITY,
			parent_id INTEGER,
			signature VARCHAR(40) NOT NULL,
			message TEXT,
			timestamp TIMESTAMP NOT NULL,
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

func (db *DBAccess) createFolderNodeTable() {
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

func (db *DBAccess) createBranchTable() {
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
				ON UPDATE CASCADE
			)`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating Branch table")
		fmt.Println(err.Error())
		panic(err)
	}
}

func (db *DBAccess) createFileNodeTable() {
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

func (db *DBAccess) createUserTable() {
	query := `CREATE TABLE IF NOT EXISTS "GardenUser" (
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

func (db *DBAccess) createRepositoryTable() {
	query := `CREATE TABLE IF NOT EXISTS Repository (
			id INTEGER PRIMARY KEY
				GENERATED ALWAYS AS IDENTITY,
			name VARCHAR(40) NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (user_id) REFERENCES "GardenUser" (id)
				ON DELETE CASCADE 
				ON UPDATE CASCADE
			)`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println("Error creating Repository table")
		fmt.Println(err.Error())
		panic(err)
	}
}

func (db *DBAccess) setupDefaultUser() int64 {
	user := types.User{
		Name:     "test",
		Password: "test",
		Email:    "test@email.com",
	}
	u, err := db.GetUserByUsername(user.Name)
	if err != nil {
		if !errors.As(err, &sql.ErrNoRows) {
			panic(fmt.Errorf("error checking for default user: %w", err))
		}
	} else {
		return u.ID
	}

	id, err := db.InsertUser(&user)
	if err != nil {
		panic(fmt.Errorf("error creating default user: %w", err))

	}
	return id
}

func (db *DBAccess) setupDefaultRepository(userId int64) int64 {
	id, err := db.InsertRepository("test", userId)
	if err != nil {
		panic(fmt.Errorf("error creating default repository: %w", err))
	}
	return id
}
