
package handlers
 import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
		"strconv"
  )

func Insertmany(db *sql.DB) http.HandlerFunc{
	return func( w http.ResponseWriter, r *http.Request){
		tx, err := db.Begin()
		if err != nil{
			log.Fatal(err)
		}
		stmt,err := tx.Prepare("INSERT INTO tasks(name,status) VALUES(?,?)")
		if err!=nil{
			log.Fatal(err)
		}
		defer stmt.Close()

		for i:=1;i<=10000;i++{
			name:= "Task"+ strconv.Itoa(i)
			status := "pending"
			if i%2==0{
				status="Done"
			}
			_,err := stmt.Exec(name,status)
			if err!=nil{
				log.Println("error in the data inserting ")
			}
		}
  
		err = tx.Commit()
		if err!=nil{
			log.Fatal(err)
		}

		json.NewEncoder(w).Encode(map[string]string{
		"message":"the data inserted in the database succesfully",
	})

	}
}