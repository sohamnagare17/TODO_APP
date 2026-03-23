package handlers

 import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
		
  )

func DeleteTask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

		id := r.URL.Query().Get("id")

		query :=`DELETE FROM tasks
		        WHERE id=?`

		_,err := db.Exec(query,id)
		
		if err!=nil{
			log.Println("error while deleting the task")
			return 
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message":"delete the task succesfully",
		})
	}
}