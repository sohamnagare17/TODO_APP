
package handlers
 import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
	
		"go-sqlite/models"
		"strconv"
  )


func ViewTask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter , r *http.Request){

		  limitstr := r.URL.Query().Get("limit")
		  afterindestr := r.URL.Query().Get("afterind")

		  limit,_:= strconv.Atoi(limitstr)
		  afterind,_ := strconv.Atoi(afterindestr)

		  rows, err := db.Query("SELECT * FROM tasks WHERE id > ? LIMIT  ?" ,afterind,limit)
		  if err!=nil{
			log.Println("error in fetching the parameters")
		  }
 defer rows.Close()

		  var list []models.Task

		  for rows.Next(){
               var task models.Task
			   err := rows.Scan(&task.ID,&task.NAME,&task.STATUS)

			   if err!=nil{
				log.Println("wrong in the scanning the data from the rows ")
			   }

			   list = append(list,task)
		  }

      json.NewEncoder(w).Encode(list)

	}


}