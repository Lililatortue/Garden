package sql

import "fmt"

func (db *DBAccess) setup() {
	db.createUserTable()
	db.createRepositoryTable()
	db.createFolderNodeTable()
	db.createGardenTagTable()
	db.createFileNodeTable()
	db.createBranchTable()
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
	query := `CREATE TABLE IF NOT EXISTS "User" (
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
			FOREIGN KEY (user_id) REFERENCES "User" (id)
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
