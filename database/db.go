package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func Dbinit() *sql.DB {

	dsn := "root:root@tcp(localhost:3307)/todo?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Error connecting to DB")
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Println("DB not reachable")
		log.Fatal(err)
	}

	// USERS TABLE
	query1 := `CREATE TABLE IF NOT EXISTS users(
	   user_id INT AUTO_INCREMENT PRIMARY KEY,
	   username VARCHAR(255),
	   email VARCHAR(255) UNIQUE
	);`

	_, err = db.Exec(query1)
	if err != nil {
		log.Println("Error creating users table")
		log.Fatal(err)
	}

	// TASKS TABLE
	query2 := `CREATE TABLE IF NOT EXISTS tasks1(
	    id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		status VARCHAR(50),
		userid INT,
		createdAt DATETIME,
		updatedAt DATETIME,
		FOREIGN KEY(userid) REFERENCES users(user_id)
	);`

	_, err = db.Exec(query2)
	if err != nil {
		log.Println("Error creating tasks table")
		log.Fatal(err)
	}

	// UNIQUE INDEX
	_, err = db.Exec(`CREATE UNIQUE INDEX singelrecords ON tasks1(name, status, userid);`)
	if err != nil {
		log.Println("Index may already exist:", err)
	}

	log.Println("✅ Tables created successfully")
	return db
}
