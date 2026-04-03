package db

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

func Dbinit() *sql.DB {
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		log.Println("here is err1")
		log.Fatal(err)
	}

	query1 := `CREATE TABLE IF NOT EXISTS users(
	   user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	   username TEXT,
	   email TEXT UNIQUE
	);`

	_, err = db.Exec(query1)

	if err != nil {
		log.Println("here is err12")
		log.Fatal(err)
	}

	query := `CREATE TABLE IF NOT EXISTS tasks1(
	    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT ,
		status TEXT,
		userid INTEGER,
		createdAt TEXT,
		updatedAt TEXT,
		FOREIGN KEY(userid) REFERENCES users(user_id)
	);`

	_, err = db.Exec(query)

	if err != nil {
		log.Println("here is err123")
		log.Fatal(err)
	}

	// db.Exec(`CREATE INDEX IF NOT EXISTS IDX ON tasks(id)`)

	db.Exec(`CREATE UNIQUE INDEX singelrecords ON tasks1(name, status, userid);`)
	db.Exec("PRAGMA foreign_keys = ON")
	log.Println("table creates succesfully")
	return db
}
