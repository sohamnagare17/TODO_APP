package handlers

import (
	"database/sql"
	"encoding/json"
	"go-sqlite/models"
	"log"
	"net/http"
	"strconv"
	"time"
)


func InsertTask(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
        var newtask models.Task;

		uid:=r.URL.Query().Get("userid")
		userid,err:=strconv.Atoi(uid)
		if err!=nil{
			log.Println("enter valid user id")
		}

		err= json.NewDecoder(r.Body).Decode(&newtask)

		if err!=nil{
		log.Println("error in fetching the data")
		}

		now:= time.Now().Format("2006-01-02 15:04:05")

		query := `INSERT INTO tasks (name,status,createdAt,updatedAt,userId) 
				VALUES(?,?,?,?,?)`

		_, err = db.Exec(query,newtask.NAME,newtask.STATUS,now,now,userid)

		if err != nil{
		log.Println("somthing went wrong to inserting the data ")
		return 
		}

		json.NewEncoder(w).Encode(map[string]string{
		"message":"the task inserted succesfully into database ",
		})

	}
}