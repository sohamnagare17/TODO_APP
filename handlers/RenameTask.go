package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func RenameTask(db *sql.DB) http.HandlerFunc{
	return func (w http.ResponseWriter , r *http.Request){

		id:= r.URL.Query().Get("id")
		name := r.URL.Query().Get("name")
		uid:=r.URL.Query().Get("userId")
		userId,err:=strconv.Atoi(uid)
		if err!=nil{
			log.Println("Enter valid userId")
		}
		taskid,err:=strconv.Atoi(id)
		if err!=nil{
			log.Println("enter valid user id")
		}

		query:=`UPDATE tasks
		        SET name=? , updatedAt=CURRENT_TIMESTAMP
				WHERE  id=? AND userId=?`

		_,err= db.Exec(query,name,taskid,userId)
		if err!=nil{
			log.Println("error in updatting the data ")
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message":"rename the task succesfully",
		})
	}
}
