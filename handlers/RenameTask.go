package handlers

  import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
		
  )



func RenameTask(db *sql.DB) http.HandlerFunc{
	return func (w http.ResponseWriter , r *http.Request){
	

		id:= r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")

		query:=`UPDATE tasks
		        SET name=?
				WHERE  id=? `

		_,err:= db.Exec(query,name,id)

		if err!=nil{
			log.Println("error in updatting the data ")
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message":"rename the task succesfully",
		})
	}
}
