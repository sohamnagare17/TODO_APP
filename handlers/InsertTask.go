package handlers

 import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
		"go-sqlite/models"
  )


func InsertTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
        var newtask models.Task;

		err := json.NewDecoder(r.Body).Decode(&newtask)

		if err!=nil{
		log.Println("error in fetching the data")
		}

		query := `INSERT INTO tasks (name , status) VALUES(?,?)`

		_, err = db.Exec(query,newtask.NAME,newtask.STATUS)

		if err != nil{
		log.Println("somthing went wrong to inserting the data ")
		return 
		}

		json.NewEncoder(w).Encode(map[string]string{
		"message":"the task inserted succesfully into database ",
		})

	}
}