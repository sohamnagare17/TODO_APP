package testutils

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func SetupTestDb() *sql.DB {

	// ✅ IMPORTANT: enable time parsing
	db, err := sql.Open("sqlite3", "file::memory:?_loc=auto&parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	createTable(db)
	insertTestData(db)

	log.Println("table created successfully")
	return db
}

func createTable(db *sql.DB) {

	queryUsers := `CREATE TABLE IF NOT EXISTS users(
	   user_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	   username TEXT,
	   email TEXT UNIQUE
	);`

	_, err := db.Exec(queryUsers)
	if err != nil {
		log.Fatal(err)
	}

	queryTasks := `CREATE TABLE IF NOT EXISTS tasks1(
	    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name TEXT,
		status TEXT,
		userid INTEGER,
		createdAt DATETIME,
		updatedAt DATETIME,
		FOREIGN KEY(userid) REFERENCES users(user_id)
	);`

	_, err = db.Exec(queryTasks)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE UNIQUE INDEX singelrecords ON tasks1(name, status, userid);`)
	if err != nil {
		log.Fatal(err)
	}

	db.Exec("PRAGMA foreign_keys = ON")
}

// func insertTestData(db *sql.DB) {

// 	now := time.Now().UTC()

// 	// insert user
// 	_, err := db.Exec(`
// 	INSERT INTO users(username, email) VALUES
// 	('testuser', 'test@example.com')
// 	`)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// insert tasks
// 	_, err = db.Exec(`
// 	INSERT INTO tasks1(name, status, userid, createdAt, updatedAt)
// 	VALUES (?, ?, ?, ?, ?),
// 	       (?, ?, ?, ?, ?),
// 	       (?, ?, ?, ?, ?)`,
// 		"Task 1", "pending", 1, now, now,
// 		"Task 2", "done", 1, now, now,
// 		"Task 3", "pending", 2, now, now,
// 	)

//		if err != nil {
//			log.Fatal(err)
//		}
//	}
func insertTestData(db *sql.DB) {

	now := time.Now().UTC()

	// ✅ Insert 2 users
	_, err := db.Exec(`
	INSERT INTO users(username, email) VALUES
	('user1', 'u1@example.com'),
	('user2', 'u2@example.com')
	`)
	if err != nil {
		log.Fatal(err)
	}

	// ✅ Now both userid 1 & 2 are valid
	_, err = db.Exec(`
	INSERT INTO tasks1(name, status, userid, createdAt, updatedAt)
	VALUES (?, ?, ?, ?, ?),
	       (?, ?, ?, ?, ?),
	       (?, ?, ?, ?, ?)`,
		"Task 1", "pending", 1, now, now,
		"Task 2", "done", 1, now, now,
		"Task 3", "pending", 2, now, now,
	)

	if err != nil {
		log.Fatal(err)
	}
}
