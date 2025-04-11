package sql

func (db *DBAccess) setup() {
	db.createGardenTagTable()
	db.createFolderNodeTable()
}

func (db *DBAccess) createGardenTagTable() {
	query := `CREATE TABLE IF NOT EXISTS GardenTag (
    			id INTEGER PRIMARY KEY,
    			signature VARCHAR(40) NOT NULL,
    			message STRING,
    			timestamp TIMESTAMP NOT NULL,
    			tree_id INTEGER NOT NULL,
    			FOREIGN KEY (tree_id) REFERENCES HashTree (id)
    				ON DELETE CASCADE 
    				ON UPDATE CASCADE,
				)`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}

func (db *DBAccess) createFolderNodeTable() {
	query := `CREATE TABLE IF NOT EXISTS FolderNode (
  				id INTEGER PRIMARY KEY,  
)`
	_, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
}
