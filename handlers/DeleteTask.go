package handlers

 import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
	"strconv"
		
  )

func DeleteTask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){

idstr := r.URL.Query().Get("id")

	 if idstr==""{
		log.Println("Id is required plz provide id")
		return
	 }

	 id,err := strconv.Atoi(idstr)
	 if err!=nil{
		log.Println("id must be number ")
		return
	 }

		query :=`DELETE FROM tasks
		        WHERE id=?`

		_,err1 := db.Exec(query,id)
		
		if err1!=nil{
			log.Println("error while deleting the task")
			return 
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message":"delete the task succesfully",
		})
	}
}