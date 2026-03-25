
package db
import(
	"database/sql"
		"log"
	_ "github.com/mattn/go-sqlite3"
)

func Dbinit() *sql.DB{
	db, err := sql.Open("sqlite3", "./test.db")
	if err!=nil{
		log.Fatal(err)
	}
	

	query := `CREATE TABLE IF NOT EXISTS tasks(
	     id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		 name TEXT ,
		 status TEXT
	
	);`

	_,err = db.Exec(query)

	if err!=nil{
		log.Fatal(err)
	}

	db.Exec(`CREATE INDEX IF NOT EXISTS IDX ON tasks(id)`)
	db.Exec(`CREATE INDEX IF NOT EXISTS IDX_NAME ON tasks(name)`)
	log.Println("table creates succesfully")
	return db
}