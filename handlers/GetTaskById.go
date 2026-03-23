package handlers
 import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
		"go-sqlite/models"
  )

func gettask(db *sql.DB) http.HandlerFunc{
	return func(w http.ResponseWriter , r *http.Request){
		var task models.Task;

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

	 query := `SELECT id,name,status FROM tasks where id=?`

	 err:= db.QueryRow(query,id).Scan(&task.ID,&task.NAME,&task.STATUS)

	 if err!=nil{
		log.Println("error in fetching the data");
		return 
	 }
	 json.NewEncoder(w).Encode(task)
	}
}