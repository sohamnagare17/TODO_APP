
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

		  limit,err:= strconv.Atoi(limitstr)
          if err!=nil{
			 log.Println("plz provide valid limit value")
		  }

		  afterind,err1:= strconv.Atoi(afterindestr)

		  if err1!=nil{
			log.Println("plz provide valid after index value")
		  }

		  rows, err := db.Query("SELECT * FROM tasks WHERE id > ? LIMIT  ?" ,afterind,limit)
		  if err!=nil{
			log.Println("error in fetching the parameters")
		  }
 defer rows.Close()

		  var list []models.Task

		  for rows.Next(){
               var task models.Task
			   err := rows.Scan(&task.Id,&task.Name,&task.Status)

			   if err!=nil{
				log.Println("wrong in the scanning the data from the rows ")
			   }

			   list = append(list,task)
		  }
          nextcursor := afterind
		  if (len(list)>0){
			nextcursor = list[len(list)-1].Id
		  }

      json.NewEncoder(w).Encode(map[string]interface{}{
		"list:":list,
		"next_cursor":nextcursor,
		"has_more":len(list)==limit,
		"afterind":afterind,
	  })

	}


}