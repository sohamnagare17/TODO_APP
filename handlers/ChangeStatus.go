  package handlers

  import (
	"encoding/json"
		"net/http"
		"log"
	"database/sql"
		"strconv"
  )

func ChangeStatus(db *sql.DB) http.HandlerFunc{
	return func( w http.ResponseWriter, r *http.Request){

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

		query := `UPDATE tasks
		          SET status=CASE
				  when status="pending" THEN "DONE"
				  ELSE "pending"
				  END
				  WHERE id=?
				  `

	   _,err1 := db.Exec(query,id)
	   
	   if err1!=nil{
		 log.Println("error in update  the status")
		 return 
	   }

	   json.NewEncoder(w).Encode(map[string]string{
		"message":"update the status succesfully ",
	   })
	}
}