  package handlers

  import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
		
  )

func ChangeStatus(db *sql.DB) http.HandlerFunc{
	return func( w http.ResponseWriter, r *http.Request){

		id := r.URL.Query().Get("id")

		query := `UPDATE tasks
		          SET status=CASE
				  when status="pending" THEN "DONE"
				  ELSE "pending"
				  END
				  WHERE id=?
				  `

	   _,err := db.Exec(query,id)
	   
	   if err!=nil{
		 log.Println("error in update  the status")
		 return 
	   }

	   json.NewEncoder(w).Encode(map[string]string{
		"message":"update the status succesfully ",
	   })
	}
}